/**
 *create by phachon@163.com
 */
var Author = {

	//注册
	register : function () {
		var element = $(".register_form");
		var username = $("input[name='username']").val();
		var password = $("input[name='password']").val();
		var email = $("input[name='email']").val();
		if(username == "") {
			return Author.error("注册失败：用户名不能为空!");
		}
		if(password == "") {
			return Author.error("注册失败：密码不能为空!");
		}
		if(email == "") {
			return Author.error("注册失败：邮箱不能为空!");
		}

		function response(result) {
			if(result.code == 0) {
				Author.error(result.message);
			}
			if(result.code == 1) {
				Author.success(result.message);
			}
			if(result.redirect.url) {
				setTimeout(function() {
					location.href = result.redirect.url;
				}, result.redirect.sleep);
			}
		}

		var options = {
			dataType: 'json',
			success: response
		};

		element.ajaxSubmit(options);
		// return Author.success("恭喜！注册成功");
	},

	//登录
	login : function () {
		var element = $(".login_form");
		var username = $("input[name='username']").val();
		var password = $("input[name='password']").val();
		if(username == "") {
			return Author.error("登录失败：用户名不能为空!");
		}
		if(password == "") {
			return Author.error("登录失败：密码不能为空!");
		}

		function response(result) {
			if(result.code == 0) {
				Author.error(result.message);
			}
			if(result.code == 1) {
				Author.success(result.message);
			}
			if(result.redirect.url) {
				setTimeout(function() {
					location.href = result.redirect.url;
				}, result.redirect.sleep);
			}
		}

		var options = {
			dataType: 'json',
			success: response
		};
		element.ajaxSubmit(options);
	},

	//success message
	success: function (message) {
		$(".success_message strong").text(message);
		$(".error_message").hide();
		$(".success_message").show();
		return true;
	},

	//error message
	error: function (message) {
		$(".error_message strong").text(message);
		$(".success_message").hide();
		$(".error_message").show();
		return false;
	}
};
