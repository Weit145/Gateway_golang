package auth

import (
	"context"
	"fmt"

	GRPCauth "github.com/Weit145/proto-repo/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api  GRPCauth.AuthClient
	conn *grpc.ClientConn
}

type AuthService interface {
	CreateUser(ctx context.Context, login, email, password, username string) (*GRPCauth.Okey, error)
	RegistrationUser(ctx context.Context, token string) (*GRPCauth.CookieResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*GRPCauth.AccessTokenResponse, error)
	Authenticate(ctx context.Context, login, password string) (*GRPCauth.CookieResponse, error)
	CurrentUser(ctx context.Context, accessToken string) (*GRPCauth.CurrentUserResponse, error)
	LogOutUser(ctx context.Context, token string) error
}

func New(addr string) (*Client, error) {
	const op = "grpc.New"

	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	grpcClient := GRPCauth.NewAuthClient(conn)

	return &Client{
		api:  grpcClient,
		conn: conn,
	}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) CreateUser(ctx context.Context, login string, email string, password string, username string) (*GRPCauth.Okey, error) {
	req := &GRPCauth.UserCreateRequest{
		Login:    login,
		Email:    email,
		Password: password,
		Username: username,
	}
	resp, err := c.api.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) RegistrationUser(ctx context.Context, token string) (*GRPCauth.CookieResponse, error) {
	req := &GRPCauth.TokenRequest{
		TokenPod: token,
	}
	resp, err := c.api.RegistrationUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*GRPCauth.AccessTokenResponse, error) {
	req := &GRPCauth.CookieRequest{
		RefreshToken: refreshToken,
	}
	resp, err := c.api.RefreshToken(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) Authenticate(ctx context.Context, login, password string) (*GRPCauth.CookieResponse, error) {
	req := &GRPCauth.UserLoginRequest{
		Login:    login,
		Password: password,
	}
	resp, err := c.api.Authenticate(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) CurrentUser(ctx context.Context, accessToken string) (*GRPCauth.CurrentUserResponse, error) {
	req := &GRPCauth.UserCurrentRequest{
		AccessToken: accessToken,
	}
	resp, err := c.api.CurrentUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) LogOutUser(ctx context.Context, token string) error {
	req := &GRPCauth.TokenRequest{
		TokenPod: token,
	}
	_, err := c.api.LogOutUser(ctx, req)
	return err
}
