layui.define(['common', 'form'], function (exports) {
    var common = layui.common,
        table = layui.table,
        $ = layui.$,
        form = layui.form;

    var obj = {
        init: function (name, url) {
            common.render(name, '#clientContainer');

            var id = $('#clientId');
            var domain = $('#clientDomain');

            var active = {
                reload: function () {
                    clientTable.reload({
                        initSort: common.initSort,
                        where: {
                            id: id.val(),
                            domain: domain.val()
                        }
                    });
                }
            };

            common.init(active);

            // 表格渲染
            var clientTable = table.render({
                elem: '#clientTable',
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
                    {field: 'name', title: '名称'},
                    {field: 'domain', title: '域名', minWidth: 200},
                    {field: 'callback', title: '回调地址', minWidth: 300},
                    {field: 'secret', title: '密钥', minWidth: 130},
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
                            var btn = "<div id=\"clientTableTool\">"
                            btn += common.template.tableTool(d);
                            btn += "</div>";

                            return btn;
                        }
                    },
                ]]
            });

            // 排序
            table.on('sort(clientTable)', function (obj) {
                clientTable.reload({
                    initSort: obj,
                    where: {
                        id: id.val(),
                        domain: domain.val(),
                        order_field: obj.field,
                        order_type: obj.type
                    }
                });
            });

            // 操作按钮
            table.on('tool(clientTable)', function (obj) {
                var data = obj.data;
                var layEvent = obj.event;

                if (layEvent === 'del') { //删除
                    layer.confirm('确定操作吗？', function (index) {
                        $.ajax({
                            type: "POST",
                            url: "/api/v1/clients/" + data.id + "/toggle-status",
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
                    $.ajax({
                        type: "GET",
                        url: "/api/v1/clients/?id=" + data.id,
                    }).done(function (data, text, jqXHR) {
                        if (!common.apiResponse(data)) {
                            return;
                        }
                        form.val("clientUpdate", data.data.items[0]);
                    });
                }
            });

            // 新增相关
            // 显示弹窗
            $('.search .create-btn').on('click', function () {
                layer.open({
                    type: 1,
                    title: ['新增客户端'],
                    shade: 0.4,
                    content: $("#clientCreate").html(),
                    success: function (layero, index) {
                        form.render();
                    }
                });
            });

            // 校验回调地址
            $('body').on('click', "#checkCallbackUrl", function () {
                let url = $('#callbackUrl').val();
                $.ajax({
                    type: "GET",
                    url: "/api/v1/clients/check-callback?url=" + url,
                }).done(function (data, text, jqXHR) {
                    if (!common.apiResponse(data, layer.alert, data.data, {icon: 5})) {
                        return;
                    }
                    layer.alert("Success", {icon: 6});
                });
            });

            // 初始化对接密钥
            $('body').on('click', "#generateSecret", function () {
                $('#clientSecret').val(Math.random().toString(36).substring(2));
            });

            // 表单提交
            form.on('submit(clientSubmit)', function (data) {
                var type = "POST";
                var url = "/api/v1/clients";
                if (data.field.id !== "") {
                    type = "PUT";
                    url = "/api/v1/clients/" + data.field.id;
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
        },
    };

    exports('client', obj);
})