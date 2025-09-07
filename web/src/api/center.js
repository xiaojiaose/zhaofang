import service from '@/utils/request'

export const searchXiaoqu = (keyword) => {
    return service({
        url: '/center/xiaoqu/list',
        method: 'post',
        data: {
            keyword: keyword,
            cityId: "1",
            page: 0,
            pageSize: 1000
        }

    })
}

export const getHouseOptions = () => { //筛选用到的筛选项
    return service({
        url: '/house/type/options',
        method: 'get'
    })
}

export const uploadFile = (data) => { //上传文件
  return service({
    url: '/center/upload',
    method: 'post',
    data: data
  })
}
