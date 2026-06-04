// Package httpmw 는 공개/관리자 API 가 공유하는 횡단 HTTP 미들웨어를 제공한다.
package httpmw

import (
	"net/http"
	"slices"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS 는 대시보드/웹앱(다른 Origin)에서의 호출을 허용하는 전역 미들웨어를 만든다.
// refine 페이지네이션을 위해 X-Total-Count 헤더를 노출하고, Authorization 요청 헤더를 허용한다.
// 전역(engine.Use)으로 등록해 프리플라이트(OPTIONS, 미매칭 라우트 포함)까지 처리한다.
func CORS(origins []string) gin.HandlerFunc {
	cfg := cors.Config{
		AllowMethods:  []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"X-Total-Count"},
	}
	// 명시적 '*' 일 때만 전체 허용. 빈 목록은 교차 출처 차단(fail-closed).
	if slices.Contains(origins, "*") {
		cfg.AllowAllOrigins = true
	} else {
		cfg.AllowOrigins = origins
	}
	return cors.New(cfg)
}
