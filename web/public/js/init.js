var MemoApp = {
    Template: {},
    Form: {}
};

(function () {
    /**
     * @param token {string}
     */
    MemoApp.InitCsrf = function (token) {
        /**
         * @param method {string}
         * @returns {boolean}
         */
        function csrfSafeMethod(method) {
            // HTTP methods that do not require CSRF protection.
            return (/^(GET|HEAD|OPTIONS|TRACE)$/.test(method));
        }

        $.ajaxSetup({
            crossDomain: false,
            beforeSend: function (xhr, settings) {
                if (!csrfSafeMethod(settings.type)) {
                    xhr.setRequestHeader("X-CSRF-Token", token);
                }
            }
        });
    };

    MemoApp.InitTimeZone = function () {
        if (document.cookie) {
            var tz = jstz.determine();
            document.cookie = "memo_time_zone=" + tz.name() + ";path=/;max-age=31104000";
        }
    };

    var alertId = 0;
    /**
     * @param message {string}
     */
    MemoApp.AddAlert = function (message) {
        alertId++;
        $(".alert-banner").append("<p id='alert-banner-message-" + alertId + "'>" + message + "</p>");
        var $alertMessage = $("#alert-banner-message-" + alertId);
        $alertMessage.hide().slideDown();
        setTimeout(function () {
            $alertMessage.slideUp(function () {
                $alertMessage.remove();
            })
        }, 5000);
    };

    var BaseURL = "/";

    /**
     * @param url {string}
     */
    MemoApp.SetBaseUrl = function (url) {
        BaseURL = url;
    };

    /**
     * @return {string}
     */
    MemoApp.GetBaseUrl = function () {
        return BaseURL;
    };

    MemoApp.utf8ByteLength = function (str) {
        // returns the byte length of an utf8 string
        var s = str.length;
        for (var i = str.length - 1; i >= 0; i--) {
            var code = str.charCodeAt(i);
            if (code > 0x7f && code <= 0x7ff) s++;
            else if (code > 0x7ff && code <= 0xffff) s += 2;
            if (code >= 0xDC00 && code <= 0xDFFF) i--; //trail surrogate
        }
        return parseInt(s);
    };

    MemoApp.GetPassword = function () {
        if (!localStorage.WalletPassword || localStorage.WalletPassword.length === 0) {
            localStorage.WalletPassword = prompt("Your wallet has become locked. Please re-enter your password to unlock your wallet.");
        }
        return localStorage.WalletPassword;
    };

    MemoApp.ReEnterPassword = function () {
        var html =
            "<div id='re-enter-password-modal'>" +
            "<div class='panel panel-default'>" +
            "<div class='panel-heading'>Re-enter Password</div>" +
            "<div class='panel-body'>" +
            "<form id='re-enter-password-form'>" +
            "<label for='re-enter-password'>Password</label>" +
            "<input type='password' id='re-enter-password' autocomplete='Password'/>" +
            "</form>" +
            "</div>" +
            "</div>" +
            "</div>";
        MemoApp.Modal(new Array(15).join(html));
    };

    /**
     * @param {string} password
     */
    MemoApp.SetPassword = function (password) {
        localStorage.WalletPassword = password;
    };

    /**
     * @param {XMLHttpRequest} xhr
     */
    MemoApp.Form.ErrorHandler = function (xhr) {
        var errorMessage =
            "Error with request (response code " + xhr.status + "):\n" +
            (xhr.responseText !== "" ? xhr.responseText + "\n" : "") +
            "If this problem persists, try refreshing the page.";
        MemoApp.AddAlert(errorMessage);
    };

    /**
     * @param {string} path
     * @return {WebSocket}
     */
    MemoApp.GetSocket = function (path) {
        var loc = window.location;
        var protocol = window.location.protocol.toLowerCase() === "https:" ? "wss" : "ws";
        var socket = new WebSocket(protocol + "://" + loc.hostname + ":" + loc.port + path);

        var heartbeatInterval = setInterval(function () {
            if (socket.readyState === socket.CLOSED) {
                clearInterval(heartbeatInterval);
                return;
            }
            var wsMessage = "heartbeat";
            socket.send(JSON.stringify(wsMessage));
        }, 15000);

        return socket;
    };

    MemoApp.Events = {};

    MemoApp.URL = {
        Index: "",
        Profile: "profile",
        LoadKey: "key/load",
        LoginSubmit: "login-submit",
        SignupSubmit: "signup-submit",
        Logout: "logout",
        MemoPost: "post",
        MemoPostAjax: "post-ajax",
        MemoPostThreadedAjax: "post-threaded-ajax",
        MemoPostMoreThreadedAjax: "post-more-threaded-ajax",
        MemoNewSubmit: "memo/new-submit",
        MemoReplySubmit: "memo/reply-submit",
        MemoFollowSubmit: "memo/follow-submit",
        MemoUnfollowSubmit: "memo/unfollow-submit",
        MemoSetNameSubmit: "memo/set-name-submit",
        MemoSetProfileSubmit: "memo/set-profile-submit",
        MemoLikeSubmit: "memo/like-submit",
        MemoWait: "memo/wait",
        MemoWaitSubmit: "memo/wait-submit",
        PollCreateSubmit: "poll/create-submit",
        PollVoteSubmit: "poll/vote-submit",
        PollVotesAjax: "poll/votes-ajax",
        ProfileSettingsSubmit: "settings-submit",
        KeyChangePasswordSubmit: "key/change-password-submit",
        KeyDeleteAccountSubmit: "key/delete-account-submit",
        TopicsSocket: "topics/socket",
        TopicsMorePosts: "topics/more-posts",
        TopicsPostAjax: "topics/post-ajax",
        TopicsFollowSubmit: "topics/follow-submit",
        TopicsCreateSubmit: "topics/create-submit",
        MemoSetProfilePicSubmit: "memo/set-profile-pic-submit",
    };
})();
