package branch

import (
    "context"
    "database/sql"
    "vesaliusm/model"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

type BranchService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewBranchService(db *sqlx.DB, ctx context.Context) *BranchService {
    return &BranchService{db: db, ctx: ctx}
}

func (s *BranchService) FindByBranchId(branchId int64) (*model.Branch, error) {
    const query = `SELECT * FROM BRANCH WHERE BRANCH_ID = :branchId`
    var b model.Branch
    err := s.db.GetContext(s.ctx, &b, query, branchId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &b, nil
}

func (s *BranchService) FindByBranchName(branchName string) (*model.Branch, error) {
    const query = `SELECT * FROM BRANCH WHERE BRANCH_NAME = :branchName`
    var b model.Branch
    err := s.db.GetContext(s.ctx, &b, query, branchName)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &b, nil
}

func (s *BranchService) FindByUrl(url string) (*model.Branch, error) {
    const query = `SELECT * FROM BRANCH WHERE URL = :url`
    var b model.Branch
    err := s.db.GetContext(s.ctx, &b, query, url)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &b, nil
}

func (s *BranchService) FirstByURLLike(urlLikeStr string) (*model.Branch, error) {
    // Use concatenation to include the wildcards in the bind value
    const query = `SELECT * FROM BRANCH WHERE URL LIKE :urlLikeStr`
    var b model.Branch
    // Add '%' wildcards to the parameter
    err := s.db.GetContext(s.ctx, &b, query, "%" + urlLikeStr + "%")
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &b, nil
}

func (s *BranchService) FindAll() ([]model.Branch, error) {
    const query = `SELECT * FROM BRANCH`
    lx := make([]model.Branch, 0)
    err := s.db.SelectContext(s.ctx, &lx, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return lx, nil
}
