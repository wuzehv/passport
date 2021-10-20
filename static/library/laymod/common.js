layui.define(function (exports) {
    var $ = layui.jquery;

    var obj = {
        // 公共配置项
        page: true,
        toolbar: false,
        height: 'full-250',
        cellMinWidth: 80,
        autoSort: false,
        successCode: 0,
        statusEnable: 1,
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
        render: function (name, id) {
            $('#breadcrumb-name').text(name);
            this.container.empty().html($(id).html());
        },
        init: function (active) {
            // 检索
            $('.search .search-btn').on('click', function () {
                var type = $(this).data('type');
                active[type] ? active[type].call(this) : '';
            });

            $('.search input').on('keyup', function (e) {
                if (e.which === 13) {
                    active["reload"].call(this);
                }
            });
        },
        template: {
            status: function(status) {
                if (status === this.statusEnable) {
                    return '启用'
                } else {
                    return '禁用'
                }
            },
            tableTool: function(data) {
                var btn = "<a class=\"layui-btn layui-btn-xs\" lay-event=\"edit\">编辑</a>";

                var c = 'layui-bg-red';
                var text = '禁用';

                if (data.status === this.statusEnable) {
                    c = 'layui-bg-green';
                    text = '启用';
                }

                btn += "<a class=\"layui-btn " + c + " layui-btn-xs\" lay-event=\"del\">" + text + "</a>";

                return btn;
            },
            formatDateTime: function(t) {
                return moment(t).format("YYYY-MM-DD HH:mm:ss");
            },
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
        },
    }

    exports('common', obj);
})