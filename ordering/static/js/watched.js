function Watched () {
  this._watched = [];
  this.enabled =  false;

  this.load()
}

Watched.prototype.save = function () {
  localStorage.setItem('watched', JSON.stringify(this));
}

Watched.prototype.load = function () {
  var watched = localStorage.getItem('watched');

  if (history) {
    try {
      var parsed = JSON.parse(watched);
      Object.assign(this, parsed)
    } catch (e) {
      console.log('error parsing', e)
    }
  }
}

Watched.prototype.enable = function () {
  this.enabled = true;
  this.save();

  this._watched.forEach(this.mark)
}

Watched.prototype.disable = function () {
  $('.episodes')

  this.enabled = false;
  this.save();

  this._watched.forEach(this.unmark)
}

Watched.prototype.toggle = function (episode) {
  var id = $(episode).data('unique');

  if (this.has(id)) {
    this.remove(id);
  } else {
    this.add(id);
  }
}

Watched.prototype.has = function (id) {
  return this._watched.includes(id);
}

Watched.prototype.add = function (id) {
  this._watched.push(id);
  this.mark(id);
  this.save();
}

Watched.prototype.remove = function (id) {
  this._watched = this._watched.filter(w => w !== id);
  this.unmark(id);
  this.save();
}

Watched.prototype.mark = function (id) {
  $("[data-unique=" + id + "]").addClass('watched');
}

Watched.prototype.unmark = function (id) {
  $("[data-unique=" + id + "]").removeClass('watched');
}
