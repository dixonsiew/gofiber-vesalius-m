package model

import (
    "math"
    "strconv"
)

type PagedList struct {
    List       any `json:"list"`
    Total      int `json:"total"`
    TotalPages int `json:"totalPages"`
}

type Pager struct {
    Total    int
    PageNum  int
    PageSize int
}

func (o *Pager) SetPageSize(pageSize int) {
    if (o.Total < pageSize || pageSize < 1) && o.Total > 0 {
        o.PageSize = o.Total
    } else {
        o.PageSize = pageSize
    }

    if o.GetTotalPages() < o.PageNum {
        o.PageNum = o.GetTotalPages()
    }

    if o.PageNum < 1 {
        o.PageNum = 1
    }
}

func (o *Pager) GetLowerBound() int {
    return (o.PageNum - 1) * o.PageSize
}

func (o *Pager) GetUpperBound() int {
    x := o.PageNum * o.PageSize
    if o.Total < x {
        x = o.Total
    }

    return x
}

func (o *Pager) GetTotalPages() int {
    v := float64(o.Total) / float64(o.PageSize)
    x := math.Ceil(v)
    return int(x)
}

func GetPager(total int, page string, limit string) Pager {
    pageNum, _ := strconv.Atoi(page)
    pageSize, _ := strconv.Atoi(limit)
    pg := Pager{
        Total:    total,
        PageNum:  pageNum,
        PageSize: pageSize,
    }
    pg.SetPageSize(pageSize)
    return pg
}
