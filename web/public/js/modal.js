(function () {
    var $modalWrapper,
        $siteCover,
        $siteWrapper,
        $innerWrapper,
        $document,
        $window,
        windowHeight,
        scrollTop;
    $(function () {
        $modalWrapper = $("#site-modal-wrapper");
        $siteCover = $("#site-wrapper-cover");
        $siteWrapper = $("#site-wrapper");
        $innerWrapper = $("#inner-site-wrapper");
        $document = $(document);
        $window = $(window);
    });
    /**
     * @param {string} title
     * @param {string} body
     */
    MemoApp.Modal = function (title, body, width) {
        var offset = (12 - width) / 2;
        var html =
            "<div class='container vertical-center'>" +
            "<div id='re-enter-password-modal' class='col-xs-12 col-sm-" + width + " col-sm-offset-" + offset + "'>" +
            "<div class='panel panel-default'>" +
            "<div class='panel-heading'>" +
            title +
            "<a data-toggle='collapse' href='#' class='close'>&times </a>" +
            "</div>" +
            "<div class='panel-body'>" +
            body +
            "</form>" +
            "</div>" +
            "</div>" +
            "</div>" +
            "</div>";
        $modalWrapper.html(html).show();
        windowHeight = $window.height();
        scrollTop = $document.scrollTop();
        $siteWrapper.css({height: windowHeight});
        $innerWrapper.css({top: -scrollTop});
        $siteWrapper.addClass("active");
        $siteCover.addClass("active");
        $document.scrollTop(0);
        $modalWrapper.find(".container").click(function (e) {
            e.stopPropagation();
        });
        $modalWrapper.find(".close").click(function(e) {
            e.preventDefault();
            MemoApp.CloseModal();
        });
        $modalWrapper.click(function () {
            MemoApp.CloseModal();
        });
    };

    MemoApp.CloseModal = function() {
        $siteWrapper.removeClass("active");
        $innerWrapper.css({top: 0});
        $document.scrollTop(scrollTop);
        $siteCover.removeClass("active");
        $modalWrapper.hide();
    };
})();
