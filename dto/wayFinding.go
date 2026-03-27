package dto

type DropdownsDto struct {
	BuildingCode string `json:"buildingCode"`
}

type BuildingsDto struct {
	BuildingCode         string `json:"buildingCode"`
	BuildingName         string `json:"buildingName" validate:"required"`
	BuildingDisplayOrder int    `json:"buildingDisplayOrder" validate:"required,numeric"`
}

type FloorsDto struct {
	FloorCode          string `json:"floorCode"`
	FloorName          string `json:"floorName" validate:"required"`
	FloorDisplayOrder  int    `json:"floorDisplayOrder" validate:"required,numeric"`
	FloorImageRaw      string `json:"floorImageRaw" validate:"required"`
}

type LocationsDto struct {
	LocationBuildingCode string `json:"locationBuildingCode" validate:"required"`
	LocationFloorCode    string `json:"locationFloorCode" validate:"required"`
	LocationTypeCode     string `json:"locationTypeCode" validate:"required"`
	LocationCode         string `json:"locationCode" validate:"required"`
	LocationName         string `json:"locationName" validate:"required"`
}

type LocationTypesDto struct {
	LocationTypeCode          string `json:"locationTypeCode"`
	LocationTypeName          string `json:"locationTypeName" validate:"required"`
	LocationTypeDisplayOrder  int    `json:"locationTypeDisplayOrder" validate:"required,numeric"`
}

type RoutesDto struct {
	RouteFromLocationId int    `json:"routeFromLocationId" validate:"required,numeric"`
	RouteToLocationId   int    `json:"routeToLocationId" validate:"required,numeric"`
	RouteFromImageRaw   string `json:"routeFromImageRaw" validate:"required"`
	RouteToImageRaw     string `json:"routeToImageRaw" validate:"required"`
}
