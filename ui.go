package main

const html = `<!DOCTYPE html>

<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>录直播 - WebUI</title>
    <link href="http://cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet">
    <script src="http://cdn.bootcss.com/jquery/3.1.1/jquery.min.js"></script>
    <script src="http://cdn.bootcss.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>
    <script src="http://cdn.bootcss.com/flv.js/1.1.0/flv.min.js"></script>
    <script language='javascript'>
        function addTaskAction() {
            $("#processing_ui").modal('show');
            $("#addTaskDialog").modal('hide');
            if ($("#addTask_pathg").attr("hidden") == "hidden") {
                addTaskCheck();
            } else
                doAddTask()
            $("#processing_ui").modal('hide');
            return false;
        }
        theTasks = null
        function showTasks() {
            $("#tasklist").html("");
            $("#processing_ui").modal('show');
            aj = $.ajax({ url: "/ajax?act=tasks", async: false });
            ret = aj.responseText;
            theTasks = JSON.parse(ret).Tasks;
            rows = "";
            row = "";
            for (i = 0; i < theTasks.length; i++) {
                v = theTasks[i];
                inf = v.LiveInfo;
                t = "[无数据]";
                if (inf != null)
                    t = inf.RoomTitle;
                s = v.Run ? "<div class=\"progress progress-striped active\"><div class=\"progress-bar progress-success\" style=\"width: 100%;\"></div></div>" : "<div class=\"progress progress-striped\"><div class=\"progress-bar progress-bar-danger\" role=\"progressbar\" style=\"width: 100%;\"></div></div>";
                m = v.M ? "循环" : "普通";
                l = !v.Run ? "未运行" : v.TimeLong
                ls = v.Live ? "<span class=\"glyphicon glyphicon-ok-circle\" />" : "<span class=\"glyphicon glyphicon-remove-circle\" />"
                buttons = "<th>"
                if (v.Run)
                    buttons += "<button onclick=\"stopBtnEvt(" + (i + 1) + ")\" class=\"btn btn-warning\" type=\"button\"><span class=\"glyphicon glyphicon-stop\" /> 停止</button>\n"
                else {
                    buttons += "<button onclick=\"startBtnEvt(" + (i + 1) + ")\" class=\"btn btn-success\" type=\"button\"><span class=\"glyphicon glyphicon-play\" /> 开始</button>\n"
                    buttons += "<button onclick=\"delBtnEvt(" + (i + 1) + ")\" class=\"btn btn-danger\" type=\"button\"><span class=\"glyphicon glyphicon-remove\" /> 删除</button>\n"
                }
                if (v.EP) {
                    buttons += "<button onclick=\"down(" + (i + 1) + ")\" class=\"btn btn-primary\" type=\"button\"><span class=\"glyphicon glyphicon-download\" /> 下载</button>\n"
                    buttons += "<button onclick=\"play(" + (i + 1) + ")\" class=\"btn btn-primary\" type=\"button\"><span class=\"glyphicon glyphicon-play-circle\" />  播放</button>\n"
                }
                buttons += "<button onclick=\"infoBtnEvt(" + (i + 1) + ")\" class=\"btn btn-info\" type=\"button\"><span class=\"glyphicon glyphicon-option-horizontal\" /> 详情</button>\n"
                buttons += "</th>"
                row = "<tr>"
                row += "<td>" + (i + 1) + "</td>"
                row += "<td><a target=\"_blank\" href=\"" + v.SiteURL + "\"  title=\"" + v.Site + "\"><img height=\"16\" width=\"16\" src=" + v.SiteIcon + " /></a></td>"
                row += "<td>" + ls + "</td>"
                row += "<td>" + m + "</td>"
                row += "<td>" + s + "</td>"
                row += "<td>" + l + "</td>"
                row += "<td>" + t + "</td>"
                row += buttons
                row += "</tr>"
                rows += row;
            }
            $("#tasklist").html(rows);
            $("#processing_ui").modal('hide');
        }

        $(document).ready(function () {
            showTasks();
        })

        function doAddTask() {
            url = $("#addTask_url").val();
            path = $("#addTask_path").val();
            if (checkPathExist(path)) {
                alert("文件(路径)已存在,请更换.")
                $("#addTaskDialog").modal('show');
                return
            }
            m = $("#addTask_m").is(':checked');
            r = $("#addTask_run").is(':checked')
            if (!m)
                path += ".flv"
            aj = $.ajax({ url: "/ajax?act=add&url=" + url + "&path=" + path + "&m=" + m + "&run=" + r, async: false });
            ret = aj.responseText;
            if (ret != "ok")
                alert("添加任务失败.");
            location.reload();
        }

        function checkPathExist(path) {
            aj = $.ajax({ url: "/ajax?act=exist&path=" + path, async: false });
            ret = aj.responseText;
            return ret == "exist";
        }

        function addTaskCheck() {
            url = $("#addTask_url").val();
            aj = $.ajax({ url: "/ajax?act=check&url=" + url, async: false });
            ret = aj.responseText;
            j = JSON.parse(ret);
            if (!j.Pass)
                alert("不支持的地址.");
            else if (j.Has) {
                if (!j.Live) {
                    $("#addTask_m").attr("checked", "checked");
                    $("#addTask_m").attr("disabled", "disabled");
                    $("#addTask_mg").attr("class", "checkbox disabled");
                }
                $("#addTask_path").val(j.Path);
                $("#addTask_url").attr("readonly", 'readonly');
                $("#addTask_pathg").removeAttr("hidden");
                $("#addTaskDialog").modal('show');
                return
            } else
                alert("不存在的房间.");
            $("#addTaskDialog").modal('show');
        }

        function startBtnEvt(o) {
            aj = $.ajax({ url: "/ajax?act=start&id=" + o, async: false });
            ret = aj.responseText;
            if (ret != "ok")
                alert("开始任务失败.");
            else
                showTasks();
        }

        function stopBtnEvt(o) {
            if (confirm("确定要停止此任务?")) {
                aj = $.ajax({ url: "/ajax?act=stop&id=" + o, async: false });
                ret = aj.responseText;
                if (ret != "ok")
                    alert("停止任务失败.");
                else
                    showTasks();
            }
        }

        function delBtnEvt(o) {
            if (confirm("确定要删除此任务?")) {
                f = confirm("删除文件(路径)?");
                aj = $.ajax({ url: "/ajax?act=del&id=" + o + "&f=" + f, async: false });
                ret = aj.responseText;
                if (ret != "ok")
                    alert("删除任务失败.");
                else
                    showTasks();
            }
        }

        function infoBtnEvt(o) {
            v = theTasks[o - 1];
            i = v.LiveInfo
            $("#info_url").val(v.URL);
            $("#info_start").val(v.Run ? v.StartTime : "未开始");
            $("#info_index").val(v.Index);
            $("#info_path").val(v.Path);
            if (v.Live) {
                $("#info_live").removeAttr('hidden');
                $("#info_nick").val(i.LiveNick);
                $("#info_d").val(i.RoomDetails);
                $("#info_i").attr("src", i.LivingIMG);
            }
            $("#info_ui").modal('show').on("");
        }

        function down(o) {
            window.location.href = "/ajax?act=get&id=" + o;
        }

        function play(o) {
            var flvPlayer = flvjs.createPlayer({
                type: 'flv',
                url: '/ajax?act=get&id=' + o
            });
            $("#player_ui").modal('show').on("hide.bs.modal", function () {
                flvPlayer.unload();
            });
            if (flvjs.isSupported()) {
                var videoElement = document.getElementById('videoElement');
                flvPlayer.attachMediaElement(videoElement);
                flvPlayer.load();
                flvPlayer.play();
            }
        }
    </script>
</head>

<body>
    <div class="container-fluid ">
        <div class="row-fluid ">
            <div class="span12 ">
                <div class="page-header ">
                    <h1><span class="glyphicon glyphicon-facetime-video"></span> 录直播
                        <small>WebUI</small>
                    </h1>
                </div>
                <h3><span class="glyphicon glyphicon-list"></span> 任务管理</h3>
                <button class="btn btn-primary" data-toggle="modal" data-target="#addTaskDialog"><span class="glyphicon glyphicon-plus"/> 添加任务...</button>
                <button class="btn btn-default" onclick="showTasks()"><span class="glyphicon glyphicon-refresh" /> 刷新列表</button>
                <table class="table ">
                    <thead>
                        <tr>
                            <th><span class="glyphicon glyphicon-minus"></span> 任务编号</th>
                            <th><span class="glyphicon glyphicon-hdd"></span> 直播平台</th>
                            <th><span class="glyphicon glyphicon-play-circle"></span> 开播状态</th>
                            <th><span class="glyphicon glyphicon-flag"></span> 任务模式</th>
                            <th><span class="glyphicon glyphicon-stats"></span> 运行状态</th>
                            <th><span class="glyphicon glyphicon-time"></span> 运行时长</th>
                            <th><span class="glyphicon glyphicon-book"></span> 房间标题</th>
                            <th><span class="glyphicon glyphicon-hand-right"></span> 其他操作</th>
                        </tr>
                    </thead>
                    <tbody id="tasklist">
                    </tbody>
                </table>
            </div>
        </div>
    </div>

    <!--添加任务遮罩层-->
    <div class="modal fade" id="addTaskDialog" tabindex="-1" role="dialog" data-backdrop="static" data-keyboard="false">
        <div class="modal-dialog">
            <div class="modal-content">
                <form onsubmit="return addTaskAction()">
                    <div class="modal-header">
                        <h4 class="modal-title"><span class="glyphicon glyphicon-plus"></span> 添加任务</h4>
                    </div>
                    <div class="modal-body">
                        <div class="form-group">
                            <label><span class="glyphicon glyphicon-film"></span> 直播地址:</label>
                            <input type="url" class="form-control" id="addTask_url" required="required" />
                        </div>
                        <div id="addTask_pathg" hidden="hidden">
                            <div class="form-group">
                                <label><span class="glyphicon glyphicon-folder-open"></span> 保存路径(文件名,自动添加后缀”.flv“):</label>
                                <input type="text" class="form-control" id="addTask_path" required="required" value="#" />
                            </div>
                            <div class="form-group">
                                <div class="checkbox" id="addTask_mg">
                                    <label><input type="checkbox" id="addTask_m" /><span class="glyphicon glyphicon-flash"></span> 循环模式</label>
                                </div>
                                <div class="checkbox">
                                    <label><input type="checkbox" id="addTask_run" checked="checked" /><span class="glyphicon glyphicon-play"></span> 立即开始</label>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-danger" onclick="location.reload()"><span class="glyphicon glyphicon-remove" /> 关闭</button>
                        <button type="submit" class="btn btn-primary"><span class="glyphicon glyphicon-ok"/> 添加</button>
                    </div>
                </form>
            </div>
        </div>
    </div>
    <!--进度条遮罩层-->
    <div class="modal fade" id="processing_ui" tabindex="-1" role="dialog" data-backdrop="static" data-keyboard="false">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <h4 class="modal-title"><span class="glyphicon glyphicon-tasks" /> 正在处理中...</h4>
                </div>
                <div class="modal-body">
                    <div class="progress progress-striped active">
                        <div class="progress-bar progress-success" style="width: 100%;"></div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <!--播放器遮罩层-->
    <div class="modal fade" id="player_ui" tabindex="-1" role="dialog" data-backdrop="static">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <h4 class="modal-title"><span class="glyphicon glyphicon-tasks" /> 播放器(按ESC快捷键退出,循环模式只能播放第一个分段)</h4>
                </div>
                <div class="embed-responsive embed-responsive-16by9">
                    <video id="videoElement" width="640" height="360" controls="controls" />
                </div>
            </div>
        </div>
    </div>
    <!--详情条遮罩层-->
    <div class="modal fade" id="info_ui" tabindex="-1" role="dialog">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <h4 class="modal-title"><span class="glyphicon glyphicon-tasks" /> 任务详情</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label><span class="glyphicon glyphicon-film"></span> 直播地址:</label>
                        <input class="form-control" id="info_url" readonly="readonly" />
                    </div>
                    <div class="form-group">
                        <label><span class="glyphicon glyphicon-time"></span> 开始时间:</label>
                        <input type="text" class="form-control" id="info_start" readonly="readonly" />
                    </div>
                    <div class="form-group">
                        <label><span class="glyphicon glyphicon-magnet"></span> 当前轮数:</label>
                        <input type="text" class="form-control" id="info_index" readonly="readonly" />
                    </div>
                    <div class="form-group">
                        <label><span class="glyphicon glyphicon-folder-open"></span> 保存路径:</label>
                        <input type="text" class="form-control" id="info_path" readonly="readonly" />
                    </div>
                    <br />
                    <div id="info_live" hidden="hidden">
                        <div class="form-group">
                            <label><span class="glyphicon glyphicon-user"></span> 主播昵称:</label>
                            <input type="text" class="form-control" id="info_nick" readonly="readonly" />
                        </div>
                        <div class="form-group">
                            <label><span class="glyphicon glyphicon-picture"></span> 直播截图:</label>
                            <br />
                            <img id="info_i" width="320" height="180" />
                        </div>
                        <div class="form-group">
                            <label><span class="glyphicon glyphicon-pencil"></span> 房间说明:</label>
                            <textarea class="form-control" id="info_d" readonly="readonly" rows="5" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</body>

</html>
`
