package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/domain/payment/model"
	"github.com/user2410/rrms-backend/internal/domain/payment/repo"
	"go.uber.org/mock/gomock"
)

func TestCreatePaymentUrl(t *testing.T) {
	testcases := []struct {
		name          string
		body          fiber.Map
		buildStubs    func(t *testing.T, repo repo.MockRepo)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			name: "OK",
			body: fiber.Map{
				"amount":    10000,
				"bank_code": "",
				"language":  "vn",
			},
			buildStubs: func(t *testing.T, repo repo.MockRepo) {
				repo.EXPECT().
					CreatePayment(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&model.PaymentModel{}, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.Equal(t, res.StatusCode, fiber.StatusCreated)

				rawData, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				var data struct {
					Url string `json:"url"`
				}
				err = json.Unmarshal(rawData, &data)
				require.NoError(t, err)

				t.Log(data)

			},
		},
	}

	for i := range testcases {
		tc := &testcases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/payment/vnpay/create_payment_url", bytes.NewReader(data))
			req.Header.Set("Content-Type", "application/json")

			srv := newTestServer(t, ctrl)
			res, err := srv.router.GetFibApp().Test(req)
			require.NoError(t, err)

			tc.checkResponse(t, res)
		})
	}
}
