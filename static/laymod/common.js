layui.define(function (exports) {
    var $ = layui.jquery;

    var obj = {
        // 公共配置项
        page: true
        , toolbar: false
        , height: 'full-250'
        , cellMinWidth: 80
        , autoSort: false
        , parseData: function (res) {
            return {
                "code": res.code,
                "msg": res.message,
                "count": res.data.total,
                "data": res.data.items
            };
        }
        , request: {
            pageName: 'page'
            , limitName: 'page_size'
        }
        , container: $('#container')
        , init: function (html) {
            this.container.empty().html(html)
        }
        , breadcrumb: function (name) {
            $('#breadcrumb-name').text(name);
        }
        , initSort: {
            field: 'id'
            , type: 'desc'
        }
        , successCode: 0
    }

    exports('common', obj);
})