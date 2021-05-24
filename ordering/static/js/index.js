/**
 * An abstraction around persistent storage for user preferences.
 * 
 * The class provides static variables with the available keys and get/put
 * functions to access the storage.
 */
class UserPreferences {
    static WATCHED_EPISODES = "watchedEpisodes";
    static HIDE_WATCHED = "hideWatched";

    constructor (){}

    /**
     * Convert the "data-series" and "data-episode-id" into a sinlge key for
     * LocalStorage.
     * 
     * @param  {Element} element the element which contains the data attributes.
     * @returns A string identifying the episode
     */
    static getEpisodeKey (element) {
        let series = element.attributes["data-series"].value;
        let episode = element.attributes["data-episode-id"].value;
        let key = `${series}-${episode}`;
        return key;
    }

    /**
     * Retrieve a value from persistent user-storage.
     * 
     * @param {String} key The name/key of the config-value to retrieve
     * @param {*} fallback The default value if the key was not yet set by the
     *   user.
     * @returns The value as defined by the user (or the default)
     */
    get(key, fallback) {
        // We want to allow "null" as fallback as well, so we have to use a
        // sentinel-value to detect any "unset" config value. To keep
        // JSON-conversion to a minimum, we keep two values
        let nullSentinelJson = '"--config--null--sentinel--"'
        let nullSentinel = "--config--null--sentinel--"
        let value = JSON.parse(
            localStorage.getItem(key) || nullSentinelJson
        );
        if (value === nullSentinel) {
            return fallback;
        }
        return value
    }

    /**
     * Sotre a new value into persistent user-config
     * 
     * @param {String} key The name/key of the config-value to store
     * @param {*} value The value to store
     */
    put(key, value) {
        localStorage.setItem(key, JSON.stringify(value))
    }

}

(function($) {
    'use strict';

    /**
     * CSS class names used to change visual display of "watched" episodes
     */
    let watchedStateClasses = {
        HIDDEN: 'hidden',
        FAINT: 'faint',
    };

    let prefs = new UserPreferences();
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
     * Updates localStorage with "watched" information of an episode
     * 
     * 
     * @param {ChangeEvent} evt The event 
     */
    var updateWatched = function (evt) {
        let newValue = evt.target.checked;
        let key = UserPreferences.getEpisodeKey(evt.target);
        let watchedEpisodes = prefs.get(UserPreferences.WATCHED_EPISODES, []);
        let index = watchedEpisodes.indexOf(key)
        if (newValue === true && index === -1) {
            watchedEpisodes.push(key);
        } else if (newValue === false && index !== -1) {
            watchedEpisodes.splice(index, 1);
        }
        prefs.put(UserPreferences.WATCHED_EPISODES, watchedEpisodes);
        setWatchedDisplayState(prefs.get(UserPreferences.HIDE_WATCHED, true));
    };

    /**
     * Hide/Show episodes, according to the "watched" state
     * 
     * @param {Boolean} doHide Whether to hide the rows or not
     */
    var setWatchedDisplayState = function (doHide) {
        let watchedEpisodes = prefs.get(UserPreferences.WATCHED_EPISODES, []);
        let watchedClass = (
            doHide ? watchedStateClasses.HIDDEN : watchedStateClasses.FAINT
        );
        $('.episode').map(function() {
            let key = UserPreferences.getEpisodeKey(this);
            for (const [_, value] of Object.entries(watchedStateClasses)) {
                $(this).removeClass(value)
            }
            if (watchedEpisodes.includes(key)) {
                $(this).addClass(watchedClass)
            } else {
                $(this).removeClass(watchedClass)
            }
        });
    };

    /**
     * Load and initialise checkbox-states for "watched" episodes.
     */
    var initWatchedStates = function () {
        let watchedEpisodes = prefs.get(UserPreferences.WATCHED_EPISODES, []);
        $('.watchedToggle').map(function() {
            let key = UserPreferences.getEpisodeKey(this);
            if (watchedEpisodes.includes(key)) {
                this.checked = true;
            } else {
                this.checked = false;
            }
        });
    }

    var registerListeners = function() {
        $('.watchedToggle').change(updateWatched);

        $('#show-watched').click(function() {
            let linkText;
            let newState;
            let currentState = prefs.get(UserPreferences.HIDE_WATCHED, true);
            if (currentState) {
                linkText = "HIDE WATCHED";
                newState = false;
            } else {
                linkText = "SHOW WATCHED";
                newState = true;
            }
            prefs.put(UserPreferences.HIDE_WATCHED, newState);
            setWatchedDisplayState(prefs.get(UserPreferences.HIDE_WATCHED, true));

            // Accessing "firstChild.innerHTML" is brittle. But I deemed this an
            // acceptable trade-off to keep code-churn minimal (unless I missed
            // something). Also, the text-value is decoupled from the HTML
            // template for the same reason. This might lead to subtle
            // display-bugs (only the text-value) if the template is updated.
            this.firstChild.innerHTML = linkText;
        });

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
        initWatchedStates();
        setWatchedDisplayState(prefs.get(UserPreferences.HIDE_WATCHED, true));

        var colourSetting = Cookies.get('colour');
        if (colourSetting === undefined) {
            Cookies.set('colour', '1');
        } else if (colourSetting === '0') {
            disableColours();
        }
    });
})(jQuery);
