const bookmarkEpisode = ({target}) => {
    const { episodeId } = target.dataset
    let bookmarks = JSON.parse(localStorage.getItem('bookmarks'))
    const index = bookmarks?.indexOf(episodeId)
    if (target.checked) {
        //add episode to localstorage
        if (!bookmarks) bookmarks = [episodeId] //bookmarks array doesnt exist. create it
        else {
            //check to see if it exists in array already
            if (index === -1) {//does not exist, push to array
                bookmarks.push(episodeId)
            }
            else return //already in array, no need to do anything
        }
    } else {
        //remove episode from localstorage
        if (!bookmarks || index === -1) return //array does not exist OR does not exist in array, no need to do anything
        else {
            bookmarks.splice(index, 1) //remove from bookmarks array
        }
    }
    if (bookmarks) localStorage.setItem('bookmarks', JSON.stringify(bookmarks)) //if bookmarks array exists, save it to local storage
}