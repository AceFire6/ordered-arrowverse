(function($) {
    'use strict';

    let localStorageKeys = {
        WATCHED_EPISODES: "watchedEpisodes",
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

    /**
     * Convert the "data-series" and "data-episode-id" into a sinlge key for
     * LocalStorage.
     * 
     * @param  {Element} element the element which contains the data attributes.
     * @returns A string identifying the episode
     */
    var getLSEpisodeKey = function (element) {
        let series = element.attributes["data-series"].value;
        let episode = element.attributes["data-episode-id"].value;
        let key = `${series}-${episode}`;
        return key;
    }

    /**
     * Updates localStorage with "watched" information of an episode
     * 
     * 
     * @param {ChangeEvent} evt The event 
     */
    var updateWatched = function (evt) {
        let newValue = evt.target.checked;
        let key = getLSEpisodeKey(evt.target);
        let watchedEpisodes = JSON.parse(
            localStorage.getItem(localStorageKeys.WATCHED_EPISODES) || "[]"
        );
        let index = watchedEpisodes.indexOf(key)
        if (newValue === true && index === -1) {
            watchedEpisodes.push(key);
        } else if (newValue === false && index !== -1) {
            watchedEpisodes.splice(index, 1);
        }
        localStorage.setItem(
            localStorageKeys.WATCHED_EPISODES, JSON.stringify(watchedEpisodes)
        );
        setWatchedDisplayState();
    };

    /**
     * Hide/Show episodes, according to the "watched" state
     */
    var setWatchedDisplayState = function () {
        let watchedEpisodes = JSON.parse(
            localStorage.getItem(localStorageKeys.WATCHED_EPISODES) || "[]"
        );
        $('.episode').map(function() {
            let key = getLSEpisodeKey(this);
            if (watchedEpisodes.includes(key)) {
                this.style.display = 'none';
            } else {
                this.style.display = 'tableRow';
            }
        });
    };

    var registerListeners = function() {
        $('.watchedToggle').change(updateWatched);

        $('#no-color').click(function() {
            if (Cookies.get('colour') === '1') {
                disableColours();
            } else {
                enableColours();
            }
        });

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

    $(document).ready(function() {
        $('#show-filter-select').select2({
          placeholder: 'Select shows to exclude...',
          allowClear: true,
          closeOnSelect: false,
          width: "100%",
        });
        registerListeners();
        setWatchedDisplayState();

        var colourSetting = Cookies.get('colour');
        if (colourSetting === undefined) {
            Cookies.set('colour', '1');
        } else if (colourSetting === '0') {
            disableColours();
        }
    });
})(jQuery);
