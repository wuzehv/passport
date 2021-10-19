layui.define('common', function (exports) {
    var common = layui.common,
        table = layui.table,
        $ = layui.$;

    var obj = {
        init: function (name, url) {
            common.render(name, '#userContainer');

            var id = $('#userId');
            var email = $('#userEmail');

            var active = {
                reload: function () {
                    userTable.reload({
                        initSort: common.initSort,
                        where: {
                            id: id.val(),
                            email: email.val()
                        }
                    });
                }
            };

            common.init(active);

            // 表格渲染
            var userTable = table.render({
                elem: '#userTable',
                url: url,
                autoSort: common.autoSort,
                page: common.page,
                toolbar: common.toolbar,
                height: common.height,
                parseData: common.parseData,
                request: common.request,
                cellMinWidth: common.cellMinWidth,
                initSort: common.initSort,
                cols: [[ //表头
                    {field: 'id', title: 'ID', sort: true, fixed: 'left'},
                    {field: 'email', title: '邮箱', minWidth: 200},
                    {field: 'realname', title: '真实姓名'},
                    {field: 'mobile', title: '手机号'},
                    {
                        field: 'status', title: '状态', templet: function (d) {
                            return common.template.status(d.status);
                        }
                    },
                    {
                        field: 'created_at', title: '创建时间', minWidth: 200, sort: true, templet: function (d) {
                            return common.template.formatDateTime(d.created_at);
                        }
                    },
                    {
                        field: '', title: '操作', minWidth: 200, templet: function (d) {
                            var btn = "<div id=\"userTableTool\">"
                            btn += common.template.tableTool(d);
                            btn += "</div>";

                            return btn;
                        }
                    },
                ]]
            });

            // 排序
            table.on('sort(userTable)', function (obj) {
                userTable.reload({
                    initSort: obj,
                    where: {
                        id: id.val(),
                        email: email.val(),
                        order_field: obj.field,
                        order_type: obj.type
                    }
                });
            });
        }
    };

    exports('user', obj);
})