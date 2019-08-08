$(function() {
    //select默认值封装
    var selects = document.getElementsByTagName("select");
    var set_default = function() {
        for (var k = 0; k < selects.length; k++) {
            var s = selects[k];
            var defaultv = s.attributes["default"] ? s.attributes["default"].value : null;
            if (defaultv) {
                for (var i = 0; i < s.length; i++) {
                    if (s[i].value == defaultv) {
                        s[i].selected = true;
                        break;
                    }
                }
            }
        }
    };
    set_default();
    window['setSelectDefault'] = set_default;
    //修正独立显示的时候,body不能滚动
    if (self == top) {
        $('html').css({ "overflow-y": 'scroll', "overflow-x": 'hidden' });
    }
    //bootstrap-confirmation 封装
    $('[data-toggle=confirmation]').each(function() {
        $(this).attr("data-link", $(this).attr('href')).removeAttr("href");
    });
    $('[data-toggle=confirmation]').confirmation({
        btnOkLabel: '确定',
        btnCancelLabel: '取消',
        onConfirm: function() {
            var $submit = $(this);
            var position = $submit.attr("data-placement")
            $.ajax({
                url: $submit.attr("data-link"),
                dataType: 'json',
                error: function(jqXHR, textStatus, textErrorThrown) {
                    $submit.notify("操作失败," + (textStatus + ':' + textErrorThrown), { position: position, className: 'error' });
                },
                success: function(data, textStatus, jqXHR) {
                    var text = data.code ? "操作成功" : "操作失败";
                    var className = data.code ? "success" : "error";
                    var msg = data.message || text
                    if (typeof msg == "object") {
                        msg = JSON.stringify(msg)
                    }
                    $submit.notify(msg, { position: position, className: className });
                    if (data.code) {
                        if (data.redirect) {
                            var time = data.redirect.sleep || 300;
                            setTimeout(function() {
                                location = data.redirect.url;
                            }, time);
                        }
                    }
                },
                type: 'GET',
            });
            return false;
        }
    });

});