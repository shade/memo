(function () {
    /**
     * @param {jQuery} $post
     */
    MemoApp.MiniProfile = function ($post) {
        var $names = $post.find(".mini-profile-name");
        var loadingHtml = "<span class='glyphicon glyphicon-refresh spinning'></span> Loading...";
        $names.each(function () {
            var $name = $(this);
            var $miniProfile = $name.find(".mini-profile");
            if (!$miniProfile) {
                return;
            }
            var profileHash = $name.attr("data-profile-hash");
            $name.hover(function () {
                $miniProfile.show();
                if ($miniProfile.data("loaded")) {
                    return;
                }
                $miniProfile.html(loadingHtml);
                $.ajax({
                    url: MemoApp.URL.ProfileMini + "/" + profileHash,
                    success: function (html) {
                        $miniProfile.html(html);
                        $miniProfile.data("loaded", true);
                    }
                })
            }, function () {
                $miniProfile.hide();
            });
        });
    };

    /**
     * @param {string} miniProfileId
     * @param {boolean} isUnfollow
     */
    MemoApp.Form.MiniProfileFollow = function (miniProfileId, isUnfollow) {
        var $nameFollow = $("#name-follow-" + miniProfileId);
        var $cancel = $nameFollow.find(".cancel");
        var $creating = $nameFollow.find(".creating");
        var $broadcasting = $nameFollow.find(".broadcasting");
        if (isUnfollow) {
            var $unfollow = $nameFollow.find(".unfollow");
            var $confirmUnfollow = $nameFollow.find(".confirm-unfollow");
            $unfollow.click(function (e) {
                e.preventDefault();
                $unfollow.addClass("hidden");
                $confirmUnfollow.removeClass("hidden");
                $cancel.removeClass("hidden");
            });
        } else {
            var $follow = $nameFollow.find(".follow");
            var $confirmFollow = $nameFollow.find(".confirm-follow");
            $follow.click(function (e) {
                e.preventDefault();
                $follow.addClass("hidden");
                $confirmFollow.removeClass("hidden");
                $cancel.removeClass("hidden");
            });
        }
        $cancel.click(function (e) {
            e.preventDefault();
            if (isUnfollow) {
                $unfollow.removeClass("hidden");
                $confirmUnfollow.addClass("hidden");
            } else {
                $follow.removeClass("hidden");
                $confirmFollow.addClass("hidden");
            }
            $cancel.addClass("hidden");
        });
    };
})();
