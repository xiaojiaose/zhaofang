import service from '@/utils/request'
// @Summary 房源审核列表
// @Produce  application/json
// @Param data body {"approvalStatus": "string","desc": true,"keyword": "string","orderKey": "string","page": 0,"pageSize": 0,"phone": "string","xiaoquId": 0}
// @Router /house/list [post]
export const getHouseList = (data) => {
  return service({
    url: '/house/list',
    method: 'post',
    data: data
  })
}

// @Summary 房源上下架
// @Produce  application/json
// @Param data body {"ids": [0],"state": 0}
// @Router /house/changeState [post]
export const changeHouseState = (data) => {
  return service({
    url: '/house/state',
    method: 'post',
    data: data
  })
}

export const createHouse = (data) => { //筛选用到的筛选项
  return service({
    url: '/house/create',
    method: 'post',
    data: data
  })
}

export const getHouseListMy = (data) => { //经纪人自己创建的房源列表
  return service({
    url: '/house/my',
    method: 'post',
    data: data
  })
}

export const batchPass = (data) => { //批量审核通过
  return service({
    url: '/house/approvalState',
    method: 'post',
    data: data
  })
}

export const deleteHouse = (id) => {
  return service({
    url: '/house/del',
    method: 'post',
    data: {
      id: id
    }
  })
}

export const editHouse = (data) => {
  return service({
    url: '/house/edit',
    method: 'post',
    data: data
  })
}



