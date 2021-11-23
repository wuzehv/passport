layui.define('common', function (exports) {
    var common = layui.common,
        table = layui.table,
        $ = layui.$;

    var obj = {
        init: function (name, url) {
            common.render(name, '#recordContainer');

            var email = $('#userEmail');
            var client = $('#clientName');

            var active = {
                reload: function () {
                    recordTable.reload({
                        initSort: common.initSort,
                        where: {
                            email: email.val(),
                            client: client.val()
                        }
                    });
                }
            };

            common.init(active);

            // 表格渲染
            var recordTable = table.render({
                elem: '#recordTable',
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
                    {field: 'client_name', title: '客户端名称'},
                    {field: 'ip_addr', title: 'IP'},
                    {field: 'user_agent', title: 'User Agent'},
                    {
                        field: 'type', title: '日志类型', templet: function (d) {
                            if (d.type === 1) {
                                return '成功';
                            } else if (d.type === 2) {
                                return '失败';
                            } else {
                                return '其他';
                            }
                        }
                    },
                    {
                        field: 'created_at', title: '创建时间', minWidth: 200, templet: function (d) {
                            return common.template.formatDateTime(d.created_at);
                        }
                    },
                ]]
            });

            // 排序
            table.on('sort(recordTable)', function (obj) {
                recordTable.reload({
                    initSort: obj,
                    where: {
                        email: email.val(),
                        client: client.val(),
                        order_field: obj.field,
                        order_type: obj.type
                    }
                });
            });
        }
    };

    exports('record', obj);
})