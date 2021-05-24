(function($) {
    'use strict';

    var openWiki = function(episode) {
        var url = $(episode).data('href');
        window.open(url, '_blank');
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

    var disableWatched = function () {
        $('#enable-watched').find('.text').text('ENABLE WATCHED')
        $('#side-buttons').addClass('hidden');
        watched.disable()
    }

    var enableWatched = function () {
        $('#enable-watched').find('.text').text('DISABLE WATCHED')
        $('#side-buttons').removeClass('hidden');
        watched.enable()
    }

    var goTo = function (place) {
        var episodes = document.querySelectorAll('.episode');
        if (place === 'top') {
            episodes[0].scrollIntoView();
        } else if (place === 'last') {
            document.querySelectorAll('.watched')[document.querySelectorAll('.watched').length-1].scrollIntoView();
        } else if (place === 'bottom') {
            episodes[episodes.length-1].scrollIntoView();
        }
    }

    var registerListeners = function() {
        $('.episode').click(function () {
            if (watched.enabled) {
                watched.toggle(this);
            } else {
                openWiki(this);
            }
        });

        $('.episode i').click(function (event) {
            event.stopPropagation();
            openWiki(this);
        })

        $('#no-color').click(function() {
            if (Cookies.get('colour') === '1') {
                disableColours();
            } else {
                enableColours();
            }
        });

        $('#enable-watched').click(function () {
            if (watched.enabled) {
                disableWatched();
            } else {
                enableWatched();
            }
        });

        $('#goto-top').click(function () { goTo('top') })
        $('#goto-last').click(function () { goTo('last') })
        $('#goto-bottom').click(function () { goTo('bottom') })

        $('.date-picker').daterangepicker({
            autoUpdateInput: false,
            showDropdowns: true,
            minDate: '2012-10-10',
            singleDatePicker: true,
            locale: {
                format: 'YYYY-MM-DD',
                cancelLabel: 'Clear'
            }
        }).on('apply.daterangepicker', function(ev, picker) {
            $(this).val(picker.startDate.format('YYYY-MM-DD'));
        });
    };

    var watched = new Watched()

    $(document).ready(function() {
        $('#show-filter-select').select2({
          placeholder: 'Select shows to exclude...',
          allowClear: true,
          closeOnSelect: false,
          width: "100%",
        });
        registerListeners();

        var colourSetting = Cookies.get('colour');
        if (colourSetting === undefined) {
            Cookies.set('colour', '1');
        } else if (colourSetting === '0') {
            disableColours();
        }

        if (watched.enabled) {
            enableWatched();
        }
    });
})(jQuery);
