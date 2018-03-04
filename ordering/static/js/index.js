(function($) {
    'use strict';

    var openWiki = function() {
        var url = $(this).parent().data('href');
        window.open(url, '_blank');
    };

    var addFilters = function() {
        var toFilterList = [];

        var selectedOptions = $('#show-filter-select').select2('data');
        $.each(selectedOptions, function() {
            toFilterList.push(this.id);
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

    var disableColours = function() {
        $('.episode, thead').addClass('no-color');
        $('#episode-list').addClass('table-striped table-hover');
        $('#no-color').find('.text').text('ENABLE COLOR');
        Cookies.set('colour', '0');
    };

    var enableColours = function() {
        $('.episode, thead').removeClass('no-color');
        $('#episode-list').removeClass('table-striped table-hover');
        $('#no-color').find('.text').text('DISABLE COLOR');
        Cookies.set('colour', '1');
    };

    var registerListeners = function() {
        // do not bind for whole row, because of action buttons
        $('.episode td:not(.watchable)').click(openWiki);
        $('#filter-button').click(addFilters);
        $('#show-hide-seen').click(function() {
            if (Cookies.get('showSeen') === '1') {
                hideSeen();
            } else {
                showSeen();
            }
        });

        $('#no-color').click(function() {
            if (Cookies.get('colour') === '1') {
                disableColours();
            } else {
                enableColours();
            }
        });
    };

    var processWatchables = function() {
        $('.watchable').each(function(e) {
            var epId = $(this).parent().attr('id');
            var i = localStorage.getItem(epId);
            if (null !== i) {
                markAsSeen(this);
            }
            else {
                markAsUnseen(this);
            }
        });
    }

    var markAsSeen = function(el) {
        var $el = $(el);
        if (!$el.hasClass('watchable')) $el = $el.parent('.watchable');
        var epId = $el.parent('.episode').attr('id');
        $el.parent('.episode').removeClass('unseen').addClass('seen');
        localStorage.setItem(epId, 'watched');
        $el.unbind('click')
           .html('<span class="fa-stack fa-lg"><i class="fa fa-eye fa-stack-1x"/><i class="fa fa-ban fa-stack-2x text-danger"/></span>')
           .bind('click', function() { markAsUnseen(this); });
    }

    var markAsUnseen = function(el) {
        var $el = $(el);
        if (!$el.hasClass('watchable')) $el = $el.parent('.watchable');
        var epId = $el.parent('.episode').attr('id');
        $el.parent('.episode').removeClass('seen').addClass('unseen');
        localStorage.removeItem(epId);
        $el.unbind('click')
           .html('<i class="fa fa-eye">')
           .bind('click', function() { markAsSeen(this); });
    }

    var hideSeen = function() {
        $('.episode').removeClass('show-seen');
        $('#show-hide-seen').find('.text').text('SHOW SEEN');
        Cookies.set('showSeen', '0');
    };

    var showSeen = function() {
        $('.episode').addClass('show-seen');
        $('#show-hide-seen').find('.text').text('HIDE SEEN');
        Cookies.set('showSeen', '1');
    };

    $(document).ready(function() {
        $('#show-filter-select').select2({
          placeholder: 'Select shows to exclude...',
          allowClear: true,
          closeOnSelect: false,
          width: "100%",
        });
        processWatchables();
        registerListeners();

        var colourSetting = Cookies.get('colour');
        if (colourSetting === undefined) {
            Cookies.set('colour', '1');
        } else if (colourSetting === '0') {
            disableColours();
        }

        var showSeenSetting = Cookies.get('showSeen');
        if (showSeenSetting === undefined || showSeenSetting === '0') {
            hideSeen();
        } else if (showSeenSetting === '1') {
            showSeen();
        }
    });
})(jQuery);
