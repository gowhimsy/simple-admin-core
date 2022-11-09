package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/suyuan32/simple-admin-core/pkg/ent"
	"github.com/suyuan32/simple-admin-core/pkg/ent/dictionary"
	"github.com/suyuan32/simple-admin-core/pkg/msg/logmsg"
	"github.com/suyuan32/simple-admin-core/pkg/statuserr"
	"github.com/suyuan32/simple-admin-core/rpc/internal/svc"
	"github.com/suyuan32/simple-admin-core/rpc/types/core"
)

type CreateOrUpdateDictionaryDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrUpdateDictionaryDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrUpdateDictionaryDetailLogic {
	return &CreateOrUpdateDictionaryDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrUpdateDictionaryDetailLogic) CreateOrUpdateDictionaryDetail(in *core.DictionaryDetail) (*core.BaseResp, error) {
	exist, err := l.svcCtx.DB.Dictionary.Query().Where(dictionary.IDEQ(in.DictionaryId)).Exist(l.ctx)

	if err != nil {
		logx.Errorw(logmsg.DatabaseError, logx.Field("detail", err.Error()))
		return nil, statuserr.NewInternalError(errorx.DatabaseError)
	}

	if !exist {
		logx.Errorw("the parent menu does not exist", logx.Field("detail", in))
		return nil, statuserr.NewInvalidArgumentError("sys.menu.parentNotExist")
	}

	if in.Id == 0 {
		err := l.svcCtx.DB.DictionaryDetail.Create().
			SetTitle(in.Title).
			SetKey(in.Key).
			SetValue(in.Value).
			SetStatus(uint8(in.Status)).
			SetDictionaryID(in.DictionaryId).
			Exec(l.ctx)

		if err != nil {
			switch {
			case ent.IsConstraintError(err):
				logx.Errorw(err.Error(), logx.Field("detail", in))
				return nil, statuserr.NewInvalidArgumentError(errorx.CreateFailed)
			default:
				logx.Errorw(logmsg.DatabaseError, logx.Field("detail", err.Error()))
				return nil, statuserr.NewInternalError(errorx.DatabaseError)
			}
		}

		return &core.BaseResp{Msg: errorx.CreateSuccess}, nil
	} else {
		err = l.svcCtx.DB.DictionaryDetail.UpdateOneID(in.Id).
			SetTitle(in.Title).
			SetKey(in.Key).
			SetValue(in.Value).
			SetStatus(uint8(in.Status)).
			Exec(l.ctx)

		if err != nil {
			switch {
			case ent.IsNotFound(err):
				logx.Errorw(err.Error(), logx.Field("detail", in))
				return nil, statuserr.NewInvalidArgumentError(errorx.TargetNotExist)
			case ent.IsConstraintError(err):
				logx.Errorw(err.Error(), logx.Field("detail", in))
				return nil, statuserr.NewInvalidArgumentError(errorx.UpdateFailed)
			default:
				logx.Errorw(logmsg.DatabaseError, logx.Field("detail", err.Error()))
				return nil, statuserr.NewInternalError(errorx.DatabaseError)
			}
		}

		return &core.BaseResp{Msg: errorx.UpdateSuccess}, nil
	}
}
