(function($) {
    'use strict';

    var openWiki = function() {
        var url = $(this).data('href');
        window.open(url, '_blank');
    };

    var addFilters = function() {
        var toFilterList = [];

        if ($('#arrow').prop('checked')) {
            toFilterList.push('arrow');
        }
        if ($('#flash').prop('checked')) {
            toFilterList.push('flash');
        }
        if ($('#legends').prop('checked')) {
            toFilterList.push('legends');
        }

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
