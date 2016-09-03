(function($) {
    'use strict';

    var openWiki = function() {
        var url = $(this).data('href');
        window.open(url, '_blank');
    };

    var addFilters = function() {
        var toFilterList = [];

        $('input[type=checkbox]').each(function() {
            if ($(this).prop('checked')) {
                toFilterList.push($(this).attr('id'));
            }
        });

        var url = '';

        if (toFilterList.length > 0) {
            url = '/hide/' + toFilterList.join('+');
        }

        url += '/';

        if (document.baseURI.match('/newest_first$')) {
            url += 'newest_first';
        }
        window.location = url;
    };

    var registerListeners = function() {
        $('.episode').click(openWiki);
        $('#filter-button').click(addFilters);
    };

    $(document).ready(registerListeners);
})(jQuery);
