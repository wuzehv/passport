layui.define('common', function (exports) {
    var common = layui.common
        , table = layui.table
        , $ = layui.$;

    var obj = {
        init: function (name, url) {
            common.breadcrumb(name);
            common.init($('#clientContainer').html());

            var id = $('#clientId');
            var domain = $('#clientDomain');

            // 表格渲染
            var clientTable = table.render({
                elem: '#clientTable'
                , url: url
                , autoSort: common.autoSort
                , page: common.page
                , toolbar: common.toolbar
                , height: common.height
                , parseData: common.parseData
                , request: common.request
                , cellMinWidth: common.cellMinWidth
                , initSort: common.initSort
                , cols: [[ //表头
                    {field: 'id', title: 'ID', sort: true, fixed: 'left'}
                    , {field: 'domain', title: '域名', minWidth: 200}
                    , {field: 'callback', title: '回调地址', minWidth: 300}
                    , {field: 'secret', title: '密钥', minWidth: 130}
                    , {field: 'status', title: '状态'}
                    , {field: 'created_at', title: '创建时间', minWidth: 200, sort: true}
                    , {field: '', title: '操作', minWidth: 200, toolbar: '#clientTableTool'}
                ]]
            });

            // 排序
            table.on('sort(clientTable)', function (obj) {
                clientTable.reload({
                    initSort: obj
                    , where: {
                        id: id.val()
                        , domain: domain.val()
                        , order_field: obj.field
                        , order_type: obj.type
                    }
                });
            });

            var active = {
                reload: function () {
                    clientTable.reload({
                        initSort: common.initSort
                        , where: {
                            id: id.val()
                            , domain: domain.val()
                        }
                    });
                }
            };

            // 检索
            $('.clientSearch .search-btn').on('click', function () {
                var type = $(this).data('type');
                active[type] ? active[type].call(this) : '';
            });

            // 操作按钮
            table.on('tool(clientTable)', function (obj) { //注：tool 是工具条事件名，test 是 table 原始容器的属性 lay-filter="对应的值"
                var data = obj.data; //获得当前行数据
                var layEvent = obj.event; //获得 lay-event 对应的值（也可以是表头的 event 参数对应的值）
                var tr = obj.tr; //获得当前行 tr 的 DOM 对象（如果有的话）

                if (layEvent === 'detail') { //查看
                    //do somehing
                    alert("1111");
                } else if (layEvent === 'del') { //删除
                    layer.confirm('真的删除行么', function (index) {
                        obj.del(); //删除对应行（tr）的DOM结构，并更新缓存
                        layer.close(index);
                        //向服务端发送删除指令
                    });
                } else if (layEvent === 'edit') { //编辑
                    //do something

                    //同步更新缓存对应的值
                    obj.update({
                        username: '123'
                        , title: 'xxx'
                    });
                } else if (layEvent === 'LAYTABLE_TIPS') {
                    layer.alert('Hi，头部工具栏扩展的右侧图标。');
                }
            });
        }
    };

    exports('record', obj);
})