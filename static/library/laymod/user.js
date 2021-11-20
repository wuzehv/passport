layui.define('common', function (exports) {
    var common = layui.common,
        table = layui.table,
        $ = layui.$,
        form = layui.form;

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

            // 操作按钮
            table.on('tool(userTable)', function (obj) {
                var data = obj.data;
                var layEvent = obj.event;

                if (layEvent === 'del') { //删除
                    layer.confirm('确定操作吗？', function (index) {
                        $.ajax({
                            type: "POST",
                            url: "/api/v1/users/" + data.id + "/toggle-status",
                        }).done(function (data, text, jqXHR) {
                            if (!common.apiResponse(data)) {
                                return;
                            }
                            layer.close(index);
                            active.reload();
                        });
                    });
                } else if (layEvent === 'edit') { //编辑
                    $('.search .create-btn').click();
                    $('input[name="email"]').disable();
                    $.ajax({
                        type: "GET",
                        url: "/api/v1/users/?id=" + data.id,
                    }).done(function (data, text, jqXHR) {
                        if (!common.apiResponse(data)) {
                            return;
                        }
                        form.val("userUpdate", data.data.items[0]);
                    });
                }
            });

            // 新增相关
            // 显示弹窗
            $('.search .create-btn').on('click', function () {
                layer.open({
                    type: 1,
                    title: ['新增用户'],
                    shade: 0.4,
                    content: $("#userCreate").html(),
                    success: function (layero, index) {
                        form.render();
                    }
                });
            });

            // 表单提交
            form.on('submit(userSubmit)', function (data) {
                var type = "POST";
                var url = "/api/v1/users";
                if (data.field.id !== "") {
                    type = "PUT";
                    url = "/api/v1/users/" + data.field.id;
                }
                $.ajax({
                    type: type,
                    data: data.field,
                    url: url,
                }).done(function (data) {
                    if (!common.apiResponse(data)) {
                        return;
                    }
                    layer.closeAll('page');
                    active.reload();
                });
                return false;
            });
        }
    };

    exports('user', obj);
})