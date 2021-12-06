//go:generate mockgen -source=handler.go -destination=mocks/mock_gateway.go -package=mocks
package transporthttp

import (
	"context"
	"encoding/json"
	amountV1 "github.com/jacktantram/payments-api/build/go/shared/amount/v1"
	paymentsV1 "github.com/jacktantram/payments-api/build/go/shared/payment/v1"
	"github.com/jacktantram/payments-api/services/payment-gateway/internal/domain"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"strings"
	"time"
)

const (
	CardNumberLen = 16
	CVVLen        = 3
	CurrencyLen   = 3
	ExpiryMonLen  = 12
)

func HandleRoutes(h Handler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/authorize", h.AuthorizeHandler).Methods(http.MethodPost)
	r.HandleFunc("/capture", h.CaptureHandler).Methods(http.MethodPost)
	r.HandleFunc("/refund", h.RefundHandler).Methods(http.MethodPost)
	r.HandleFunc("/void", h.VoidHandler).Methods(http.MethodPost)

	return r
}

type Gateway interface {
	CreatePayment(ctx context.Context, amount *amountV1.Money, method domain.PaymentMethod) (*paymentsV1.Payment, error)
	Capture(ctx context.Context, paymentID string, amount uint64) (*paymentsV1.Payment, error)
	Refund(ctx context.Context, paymentID string, amount uint64) (*paymentsV1.Payment, error)
	Void(ctx context.Context, paymentID string) (*paymentsV1.Payment, error)
}

type Handler struct {
	gateway Gateway
}

func NewHandler(processorClient Gateway) (Handler, error) {
	if processorClient == nil {
		return Handler{}, errors.New("gateway client is nil")
	}
	return Handler{gateway: processorClient}, nil
}

func (h Handler) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		http.Error(w, "no body supplied", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var authorizationRequest CreateAuthorizationRequest
	if err := json.NewDecoder(r.Body).Decode(&authorizationRequest); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	validateRequest := func() error {
		if authorizationRequest.Amount == nil {
			return errors.New("invalid amount: cannot be missing")
		}
		if authorizationRequest.Amount.MinorUnits == 0 {
			return errors.New("invalid amount.minor_units: cannot be zero")
		}
		if len(authorizationRequest.Amount.Currency) != CurrencyLen {
			return errors.Errorf("invalid amount.currency: must be length of %d", CurrencyLen)
		}
		if authorizationRequest.Card == nil {
			return errors.New("missing payment method: cannot be empty")
		}
		if len(authorizationRequest.Card.CardNumber) != CardNumberLen && len(strings.ReplaceAll(authorizationRequest.Card.CardNumber, " ", "")) != CardNumberLen {
			return errors.Errorf("invalid payment_method.card.card_number: length not equal to %d", CardNumberLen)
		}
		if len(authorizationRequest.Card.Cvv) != CVVLen {
			return errors.Errorf("invalid payment_method.card.cvv: length not equal to %d", 3)
		}
		if authorizationRequest.Card.Expiry == nil {
			return errors.New("missing payment_method.card.expiry: cannot be empty")
		}
		if authorizationRequest.Card.Expiry.Month > 12 {
			return errors.Errorf("invalid payment_method.card.expiry.month: expiry month cannot exceed %d", ExpiryMonLen)
		}
		if int(authorizationRequest.Card.Expiry.Year) < time.Now().Year() {
			return errors.New("missing payment_method.card.expiry.year: cannot be in the past")
		}
		return nil
	}
	if err := validateRequest(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	logFields := log.Fields{
		"amount.minor_units": authorizationRequest.Amount.MinorUnits,
		"amount.currency":    authorizationRequest.Amount.Currency,
		// Ideally a card PaymentID/token would be better
		"card.first_six": authorizationRequest.Card.CardNumber[:6],
		"card.last_four": authorizationRequest.Card.CardNumber[len(authorizationRequest.Card.CardNumber)-3:],
		// could probably be injected via a middleware
		"url": "/authorize",
	}

	fn := func() error {
		paymentResponse, err := h.gateway.CreatePayment(r.Context(), authorizationRequest.Amount, domain.PaymentMethod{Card: authorizationRequest.Card})
		if err != nil {
			return err
		}
		paymentBytes, err := protojson.Marshal(paymentResponse)
		if err != nil {
			logFields["payment.id"] = paymentResponse.Id
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(paymentBytes)
		if err != nil {
			return err
		}
		return nil
	}

	if err := fn(); err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error("failed to process authorization request")
		http.Error(w, "Oops something went wrong", http.StatusInternalServerError)
		return
	}
}

func (h Handler) CaptureHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		http.Error(w, "no body supplied", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var captureRequest CreateCaptureRequest
	if err := json.NewDecoder(r.Body).Decode(&captureRequest); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	validateRequest := func() error {
		if captureRequest.PaymentID == "" {
			return errors.New("invalid payment_id: cannot be empty")
		}
		if captureRequest.Amount == 0 {
			return errors.New("invalid amount: cannot be zero")
		}
		return nil
	}
	if err := validateRequest(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	logFields := log.Fields{
		"payment.id": captureRequest.PaymentID,
		"amount":     captureRequest.Amount,
		"url":        "/capture",
	}

	fn := func() error {
		captureResponse, err := h.gateway.Capture(r.Context(), captureRequest.PaymentID, captureRequest.Amount)
		if err != nil {
			return err
		}

		paymentBytes, err := protojson.Marshal(captureResponse)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(paymentBytes)
		if err != nil {
			return err
		}
		return nil
	}

	if err := fn(); err != nil {
		if errors.Is(err, domain.ErrNoPayment) {
			http.Error(w, "payment not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrNotPermitted) {
			http.Error(w, "capture not allowed", http.StatusForbidden)
			return
		}
		logFields["error"] = err
		log.WithFields(logFields).Error("failed to process capture request")
		http.Error(w, "Oops something went wrong", http.StatusInternalServerError)
		return
	}
}

func (h Handler) RefundHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		http.Error(w, "no body supplied", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var refundRequest CreateRefundRequest
	if err := json.NewDecoder(r.Body).Decode(&refundRequest); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	validateRequest := func() error {
		if refundRequest.PaymentID == "" {
			return errors.New("invalid payment_id: cannot be empty")
		}
		if refundRequest.Amount == 0 {
			return errors.New("invalid amount: cannot be zero")
		}
		return nil
	}
	if err := validateRequest(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	logFields := log.Fields{
		"payment.id": refundRequest.PaymentID,
		"amount":     refundRequest.Amount,
		"url":        "/refund",
	}

	fn := func() error {
		refundResponse, err := h.gateway.Refund(r.Context(), refundRequest.PaymentID, refundRequest.Amount)
		if err != nil {
			return err
		}
		paymentBytes, err := protojson.Marshal(refundResponse)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(paymentBytes)
		if err != nil {
			return err
		}
		return nil
	}

	if err := fn(); err != nil {
		if errors.Is(err, domain.ErrNoPayment) {
			http.Error(w, "payment not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrNotPermitted) {
			http.Error(w, "refund not allowed", http.StatusForbidden)
			return
		}
		logFields["error"] = err
		log.WithFields(logFields).Error("failed to process refund request")
		http.Error(w, "Oops something went wrong", http.StatusInternalServerError)
		return
	}
}

func (h Handler) VoidHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		http.Error(w, "no body supplied", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var voidRequest CreateVoidRequest
	if err := json.NewDecoder(r.Body).Decode(&voidRequest); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	validateRequest := func() error {
		if voidRequest.PaymentID == "" {
			return errors.New("invalid payment_id: cannot be empty")
		}
		return nil
	}
	if err := validateRequest(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	logFields := log.Fields{
		"payment.id": voidRequest.PaymentID,
		"url":        "/void",
	}

	fn := func() error {
		voidResponse, err := h.gateway.Void(r.Context(), voidRequest.PaymentID)
		if err != nil {
			return err
		}
		paymentBytes, err := protojson.Marshal(voidResponse)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(paymentBytes)
		if err != nil {
			return err
		}
		return nil
	}

	if err := fn(); err != nil {
		if errors.Is(err, domain.ErrNoPayment) {
			http.Error(w, "payment not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrNotPermitted) {
			http.Error(w, "void not allowed", http.StatusForbidden)
			return
		}
		logFields["error"] = err
		log.WithFields(logFields).Error("failed to process void request")
		http.Error(w, "Oops something went wrong", http.StatusInternalServerError)
		return
	}
}
