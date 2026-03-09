package userPackage

type PackageCheckResult struct {
    PackageID           int64 `json:"package_id"`
    Expired             int   `json:"expired"`
    Soldout             int   `json:"soldout"`
    ExceedPurchase      int   `json:"exceedPurchase"`
    RecommendedQuantity int   `json:"recommendedQuantity"`
}
