$(function() {
    //ajax表单封装
    $("form.ajaxform").submit(function() {
        var ajaxPre = window[$(this).attr("ajaxPre")] || function() {
            return true;
        };
        var ajaxSuccess = window[$(this).attr("ajaxSuccess")] || function() {};
        var ajaxError = window[$(this).attr("ajaxError")] || function() {};
        var ajaxAlways = window[$(this).attr("ajaxAlways")] || function() {};
        var isRefresh = $(this).attr("isRefresh") != 'false';
        var $submit = $(this).find('[type="submit"]');
        $(this).ajaxSubmit({
            dataType: 'json',
            beforeSubmit: function(arr, $form, options) {
                return ajaxPre(arr, $form, options);
            },
            error: function(jqXHR, textStatus, textErrorThrown) {
                ajaxError(jqXHR, textStatus, textErrorThrown);
                $submit.notify("操作失败," + (textStatus + ':' + textErrorThrown), { position: "right", className: 'error' });
            },
            success: function(data, textStatus, jqXHR) {
                ajaxSuccess(data, textStatus, jqXHR);
                var text = data.code ? "操作成功" : "操作失败";
                var className = data.code ? "success" : "error";
                var msg = ""
                if (typeof data.message == "object") {
                    for (key in data.message) {
                        if (data.message.hasOwnProperty(key)) {
                            msg += key + data.message[key] + "\n";
                        }
                    }
                } else {
                    msg = data.message
                }
                $submit.notify(msg || text, { position: "right", className: className });
                if (data.code) {
                    if (isRefresh && data.redirect) {
                        var time = data.redirect.sleep || 300;
                        setTimeout(function() {
                            location = data.redirect.url;
                        }, time);
                    }
                }
            },
            complete: function(jqXHR, textStatus) {
                ajaxAlways(jqXHR, textStatus);
            },
            type: 'POST',
        });
        return false
    });
    //图标选择器
     
    $(function() {
        try{
            $(".icon-selector").iconSelector({
                input: '.icon'
            });
            $('.icon-placeholder').bind('icon:inserted', function(e) {
                $('.icon-placeholder-preview').html('<i class="' + this.value + '"></i>');
            });
        }catch(e){} 
    });
   
});