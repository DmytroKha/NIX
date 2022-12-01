package controllers

import (
	"NIX/internal/app"
	"NIX/internal/infra/http/requests"
	"NIX/internal/infra/http/resources"
	"fmt"
	"github.com/goccy/go-json"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	authService app.AuthService
	userService app.UserService
}

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	oauthStateString = strconv.Itoa(rand.Intn(1000000))
)

type content struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Verified bool   `json:"verified_email"`
	Picture  string `json:"picture"`
}

func NewAuthController(as app.AuthService, us app.UserService) AuthController {
	return AuthController{
		authService: as,
		userService: us,
	}
}

// NewUser godoc
// @Summary      Create a new user
// @Description  register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input   body      requests.UserRequest  true  "User body"
// @Success      201  {object}  resources.AuthDto
// @Failure      400  {string}  echo.HTTPError
// @Failure      422  {string}  echo.HTTPError
// @Failure      500  {string}  echo.HTTPError
// @Router       /auth/register [post]
func (c AuthController) Register(ctx echo.Context) error {
	var usr requests.UserRequest
	err := ctx.Bind(&usr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = ctx.Validate(&usr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	u, err := usr.ToDomainModel()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	user, token, err := c.authService.Register(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var authDto resources.AuthDto

	return ctx.JSON(http.StatusCreated, authDto.DomainToDto(token, user))

}

// LogInUser godoc
// @Summary      Log in user
// @Description  log in user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input   body      requests.UserRequest  true  "User body"
// @Success      200  {object}  resources.AuthDto
// @Failure      400  {string}  echo.HTTPError
// @Failure      422  {string}  echo.HTTPError
// @Failure      500  {string}  echo.HTTPError
// @Router       /auth/login [post]
func (c AuthController) Login(ctx echo.Context) error {
	var usr requests.UserRequest
	err := ctx.Bind(&usr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = ctx.Validate(&usr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	u, err := usr.ToDomainModel()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	user, token, err := c.authService.Login(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var authDto resources.AuthDto

	return ctx.JSON(http.StatusOK, authDto.DomainToDto(token, user))
}

// LogInUserWithGooglePartOne godoc
// @Summary      Log in user with Google (1)
// @Description  log in user with Google (1)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      201  {object}  resources.GoogleUrlDto
// @Failure      400  {string}  echo.HTTPError
// @Failure      422  {string}  echo.HTTPError
// @Failure      500  {string}  echo.HTTPError
// @Router       /auth/loginGoogle [post]
func (c AuthController) LoginGoogle(ctx echo.Context) error {

	/*err := redirectGoogleUrl(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var googleUrlDto resources.GoogleUrlDto


	return ctx.JSON(http.StatusCreated, googleUrlDto.DomainToDto(url))
	*/
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	err := ctx.Redirect(http.StatusTemporaryRedirect, url)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return nil
	//var googleUrlDto resources.GoogleUrlDto

	//return ctx.JSON(http.StatusCreated, googleUrlDto.DomainToDto(url))
}

func redirectGoogleUrl(ctx echo.Context) error {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	//err := ctx.Redirect(http.StatusTemporaryRedirect, url)
	//if err != nil {

	//	return echo.NewHTTPError(http.StatusBadRequest, err)
	//}
	return ctx.Redirect(http.StatusTemporaryRedirect, url)
}
func (c AuthController) Callback(ctx echo.Context) error {

	if ctx.FormValue("state") != oauthStateString {
		return fmt.Errorf("invalid oauth state")
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, ctx.FormValue("code"))
	if err != nil {
		//return fmt.Errorf("code exchange failed: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var authDto resources.AuthDto

	cont, err := getUserInfo(token)
	if err != nil {
		//fmt.Println(err.Error())
		//ctx.Redirect(http.StatusTemporaryRedirect, "/")
		//return err
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	var cnt content
	err = json.Unmarshal(cont, &cnt)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if !cnt.Verified {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("email is not verified"))
	}

	user, jwtToken, err := c.authService.LoginGoogle(cnt.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusCreated, authDto.DomainToDto(jwtToken, user))

}

func getUserInfo(token *oauth2.Token) ([]byte, error) {
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, nil
}
