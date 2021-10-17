layui.define(function (exports) {
    var $ = layui.jquery;

    var obj = {
        // 公共配置项
        page: true,
        toolbar: false,
        height: 'full-250',
        cellMinWidth: 80,
        autoSort: false,
        parseData: function (res) {
            return {
                "code": res.code,
                "msg": res.message,
                "count": res.data.total,
                "data": res.data.items
            };
        },
        request: {
            pageName: 'page',
            limitName: 'page_size'
        },
        container: $('#container'),
        initSort: {
            field: 'id',
            type: 'desc'
        },
        successCode: 0,
        init: function (html) {
            this.container.empty().html(html)
        },
        breadcrumb: function (name) {
            $('#breadcrumb-name').text(name);
        },
        template: {
            status: function(status) {
                if (status === 1) {
                    return '启用'
                } else {
                    return '禁用'
                }
            },
        },
        formatDateTime: function(t) {
            return moment(t).format("YYYY-MM-DD HH:mm:ss");
        },
        apiResponse: function (data, f, p1, p2) {
            if (data.code !== this.successCode) {
                if (f === undefined) {
                    layer.msg(data.message);
                } else {
                    f(p1, p2);
                }
                return false;
            }

            return true;
        }
    }

    exports('common', obj);
})