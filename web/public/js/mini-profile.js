(function () {
    var loadingHtml = "<span class='glyphicon glyphicon-refresh spinning'></span> Loading...";
    /**
     * @param {jQuery} $post
     */
    MemoApp.MiniProfile = function ($post) {
        var $names = $post.find(".mini-profile-name");
        $names.each(function () {
            var $name = $(this);
            var $miniProfile = $name.find(".mini-profile");
            if (!$miniProfile) {
                return;
            }
            $name.find(".profile-link").click(function(e) {
                e.preventDefault();
            });
            var address = $name.attr("data-profile-hash");
            $name.hover(function () {
                $miniProfile.show();
                if ($miniProfile.data("loaded")) {
                    return;
                }
                loadMiniProfile($miniProfile, address);
            }, function () {
                $miniProfile.hide();
            });
        });
    };

    /**
     * @param {jQuery} $miniProfile
     * @param {string} address
     */
    function loadMiniProfile($miniProfile, address) {
        $miniProfile.html(loadingHtml);
        $.ajax({
            url: MemoApp.URL.ProfileMini + "/" + address,
            success: function (html) {
                $miniProfile.html(html);
                $miniProfile.data("loaded", true);
            }
        })
    }

    /**
     * @param {string} miniProfileId
     * @param {string} address
     * @param {boolean} isUnfollow
     */
    MemoApp.Form.MiniProfileFollow = function (miniProfileId, address, isUnfollow) {
        var $nameFollow = $("#name-follow-" + miniProfileId);
        var $cancel = $nameFollow.find(".cancel");
        var $creating = $nameFollow.find(".creating");
        var $broadcasting = $nameFollow.find(".broadcasting");
        var submitting;
        var $miniProfile = $nameFollow.parents(".mini-profile");
        if (isUnfollow) {
            var $unfollow = $nameFollow.find(".unfollow");
            var $confirmUnfollow = $nameFollow.find(".confirm-unfollow");
            $unfollow.click(function (e) {
                e.preventDefault();
                $unfollow.addClass("hidden");
                $confirmUnfollow.removeClass("hidden");
                $cancel.removeClass("hidden");
            });
            $confirmUnfollow.click(function (e) {
                e.preventDefault();
                followSubmit(true);
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
            $confirmFollow.click(function (e) {
                e.preventDefault();
                followSubmit(false);
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

        /**
         * @param {boolean} isUnfollow
         */
        function followSubmit(isUnfollow) {
            var password = MemoApp.GetPassword();
            if (!password.length) {
                MemoApp.AddAlert("Password not set. Please re-enter and submit again.");
                MemoApp.ReEnterPassword(function() {
                    followSubmit(isUnfollow);
                });
                return;
            }

            $creating.removeClass("hidden");
            $cancel.addClass("hidden");
            if (isUnfollow) {
                $confirmUnfollow.addClass("hidden");
            } else {
                $confirmFollow.addClass("hidden");
            }
            var url;
            if (isUnfollow) {
                url = MemoApp.URL.MemoUnfollowSubmit;
            } else {
                url = MemoApp.URL.MemoFollowSubmit;
            }
            submitting = true;
            $.ajax({
                type: "POST",
                url: MemoApp.GetBaseUrl() + url,
                data: {
                    address: address,
                    password: password
                },
                success: function (followTxHash) {
                    submitting = false;
                    if (!followTxHash || followTxHash.length === 0) {
                        MemoApp.AddAlert("Server error. Please try refreshing the page.");
                        return
                    }
                    $creating.addClass("hidden");
                    $broadcasting.removeClass("hidden");
                    $.ajax({
                        type: "POST",
                        url: MemoApp.GetBaseUrl() + MemoApp.URL.MemoWaitSubmit,
                        data: {
                            txHash: followTxHash
                        },
                        success: function () {
                            submitting = false;
                            loadMiniProfile($miniProfile, address);
                        },
                        error: function () {
                            submitting = false;
                            $broadcasting.addClass("hidden");
                            $cancel.removeClass("hidden");
                            if (isUnfollow) {
                                $confirmUnfollow.removeClass("hidden");
                            } else {
                                $confirmFollow.removeClass("hidden");
                            }
                            MemoApp.Alert("Error waiting for transaction to broadcast.");
                        }
                    });
                },
                error: function (xhr) {
                    submitting = false;
                    $creating.addClass("hidden");
                    $cancel.removeClass("hidden");
                    if (isUnfollow) {
                        $confirmUnfollow.removeClass("hidden");
                    } else {
                        $confirmFollow.removeClass("hidden");
                    }
                    if (xhr.status === 401) {
                        MemoApp.AddAlert("Error unlocking key. " +
                            "Please verify your password is correct. " +
                            "If this problem persists, please try refreshing the page.");
                        MemoApp.ReEnterPassword(function () {
                            followSubmit(isUnfollow);
                        });
                        return;
                    } else if (xhr.status === 402) {
                        MemoApp.AddAlert("Please make sure your account has enough funds.");
                        return;
                    }
                    var errorMessage =
                        "Error with request (response code " + xhr.status + "):\n" +
                        (xhr.responseText !== "" ? xhr.responseText + "\n" : "") +
                        "If this problem persists, try refreshing the page.";
                    MemoApp.AddAlert(errorMessage);
                }
            });
        }
    };
})();
