const base62 = require('./base62');

class Device {
	constructor(ident, controller) {
		this.ident = ident;
		this.deviceState = {
			track: []
		}
		controller.on('update', (update) => this.handleUpdate_(update));
	}

	handleUpdate_(deviceUpdate) {
		if (deviceUpdate.ident != this.ident) {
			return;
		}

		// Volume update
		if (deviceUpdate.typ == 27) {
			this.deviceState.volume = deviceUpdate.volume
		} else if (deviceUpdate.typ == 10) { //device notify
			Object.assign(this.deviceState, deviceUpdate)
		}
	}

	convertToTracks_(ids, queued = false) {
		ids = ids.map((id) => new Uint8Array(base62.toBytes(id)))
				 .map(id => ({gid: id, queued: queued}));	
	}

	integrateLoadTracks_(ids) {
		if (ids.length == 0) {
			return [];
		}

		ids = this.convertToTracks_(ids);

		return [ids[0]]
				.concat(this.deviceState.track.filter(t => t.queued))
				.concat(ids.slice(1));
	}

	integrateQueueTracks_(ids) {
		if (ids.length == 0) {
			return [];
		}
		const tracks = this.deviceState.track;
		ids = this.convertToTracks_(ids, true);
		return 	[tracks[0]]
				.concat(tracks.slice(1).filter(t => t.queued))
				.concat(ids);
	}

	getAlbumTracks(album) {
		return album.Disc.reduce((prev, d) => prev.concat(d.track.map((t) => t.Gid)), [])
							.map((id) => base62.fromByte(id));
	}

	getUriTracks(uri) {
		const parts = uri.split(':');
		const type = parts[1];
		const id62 = parts[2];
		const id16 = Base62.toHex(id62);


		if (type == 'track') {
			return [id];
		} else if (type == 'artist') {

		} else if (type == 'album') {
			return this.controller.getAlbum(id16)
					.then((album) => this.getAlbumTracks(album));
		}
	}

	queueUri(uri) {
		const ids = this.getUriTracks(uri);
		return this.controller.replaceTracks(this.ident, this.integrateQueueTracks_(ids));
	}

	playUri(uri) {
		const ids = this.getUriTracks(uri);
		return this.controller.loadTracks(this.ident, this.integrateQueueTracks_(ids));
	}
}

module.exports = Device;
