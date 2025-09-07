import service from '@/utils/request'

export const getBuildings = (xiaoquId) => { //获取楼栋号
    return service({
        url: '/base/building',
        method: 'get',
        params: {
            id: xiaoquId
        }
    })
}

export const getUnits = (buildingId) => { //获取单元号
    return service({
        url: '/base/unit',
        method: 'get',
        params: {
            id: buildingId
        }
    })
}

export const getHouse = (unitId) => { //获取门牌号
    return service({
        url: '/base/house',
        method: 'get',
        params: {
            id: unitId
        }
    })
}


