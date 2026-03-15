package assignBranch

import (
	"context"
	"vesaliusm/database"
	"vesaliusm/model"
	"vesaliusm/service/branch"
	"vesaliusm/utils"

	"github.com/jmoiron/sqlx"
)

var AssignBranchSvc *AssignBranchService = NewAssignBranchService(database.GetDb(), database.GetCtx(), branch.BranchSvc)

type AssignBranchService struct {
    db            *sqlx.DB
    ctx           context.Context
    branchService *branch.BranchService
}

func NewAssignBranchService(db *sqlx.DB, ctx context.Context, branchService *branch.BranchService) *AssignBranchService {
    return &AssignBranchService{db: db, ctx: ctx, branchService: branchService}
}

func (s *AssignBranchService) FindAllPrimary(userId int64) ([]model.AssignBranch, error) {
    const query = `SELECT * FROM ASSIGN_BRANCH WHERE USER_ID = :userId ORDER BY ASSIGN_BRANCH_ID ASC`
    list := make([]model.AssignBranch, 0)
    err := s.db.SelectContext(s.ctx, &list, query, userId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *AssignBranchService) FindAllByUserId(userId int64) ([]model.AssignBranch, error) {
    list, err := s.FindAllPrimary(userId)
    if err != nil {
        return nil, err
    }

    for i := range list {
        b, err := s.branchService.FindByBranchId(list[i].BranchID.Int64)
        if err != nil {
            return nil, err
        }
        b.Passcode.String = ""
        b.Url.String = ""
        list[i].Branch = b
    }
    return list, nil
}

func (s *AssignBranchService) FindAllByAdminId(adminId int64) ([]model.AssignBranch, error) {
    const query = `SELECT * FROM ASSIGN_BRANCH WHERE ADMIN_ID = :adminId ORDER BY ASSIGN_BRANCH_ID ASC`
    list := make([]model.AssignBranch, 0)
    err := s.db.SelectContext(s.ctx, &list, query, adminId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *AssignBranchService) DeleteByUserIdNBranchIdNPRN(userId int64, branchId int64, prn string) error {
    const query = `DELETE FROM ASSIGN_BRANCH WHERE PRN = :prn AND USER_ID = :userId AND BRANCH_ID = :branchId`
    _, err := s.db.ExecContext(s.ctx, query, prn, userId, branchId)
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *AssignBranchService) ExistsByPRNinAssignBranch(prn string) (bool, error) {
    const query = `SELECT COUNT(*) AS COUNT FROM DUAL WHERE EXISTS (SELECT 1 FROM ASSIGN_BRANCH WHERE PRN = :prn)`
    var count int
    err := s.db.GetContext(s.ctx, &count, query, prn)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count == 1, err
}

func (s *AssignBranchService) ExistsByUserIdNBranchIdNPRNinAssignBranch(userId int64, branchId int64, prn string) (bool, error) {
    const query = `SELECT COUNT(*) AS COUNT FROM DUAL WHERE EXISTS (SELECT 1 FROM ASSIGN_BRANCH WHERE USER_ID = :userId AND BRANCH_ID = :branchId AND PRN = :prn)`
    var count int
    err := s.db.GetContext(s.ctx, &count, query, prn)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count == 1, err
}
