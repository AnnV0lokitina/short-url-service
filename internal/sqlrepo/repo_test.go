package sqlrepo

import (
	"context"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRepo_PingBD(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := pgxmock.NewConn()
	require.NoError(t, err)
	defer conn.Close(ctx)

	type input struct {
		conn PgxIface
		ctx  context.Context
	}
	tests := []struct {
		name  string
		input input
		want  bool
	}{
		{
			name: "ping positive",
			input: input{
				ctx:  ctx,
				conn: conn,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				conn: tt.input.conn,
			}
			assert.Equalf(t, tt.want, r.PingBD(tt.input.ctx), "PingBD(%v)", tt.input.ctx)
		})
	}
}
