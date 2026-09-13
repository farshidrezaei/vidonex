export namespace db {
	
	export class MediaAssetRecord {
	    id: string;
	    project_id: string;
	    file_name: string;
	    file_path: string;
	    file_type: string;
	    file_size_bytes: number;
	    duration_seconds: number;
	    width: number;
	    height: number;
	    frame_rate: number;
	    sample_rate: number;
	    channels: number;
	    thumbnail_path: string;
	    waveform_data: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new MediaAssetRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.project_id = source["project_id"];
	        this.file_name = source["file_name"];
	        this.file_path = source["file_path"];
	        this.file_type = source["file_type"];
	        this.file_size_bytes = source["file_size_bytes"];
	        this.duration_seconds = source["duration_seconds"];
	        this.width = source["width"];
	        this.height = source["height"];
	        this.frame_rate = source["frame_rate"];
	        this.sample_rate = source["sample_rate"];
	        this.channels = source["channels"];
	        this.thumbnail_path = source["thumbnail_path"];
	        this.waveform_data = source["waveform_data"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace desktop {
	
	export class ImportedAssetResult {
	    asset: db.MediaAssetRecord;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new ImportedAssetResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.asset = this.convertValues(source["asset"], db.MediaAssetRecord);
	        this.error = source["error"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SystemCapabilities {
	    os: string;
	    architecture: string;
	    cpu_count: number;
	    available_accelerators: string[];
	    ffmpeg_version: string;
	    ffprobe_version: string;
	    server_address: string;
	    server_port: number;
	
	    static createFrom(source: any = {}) {
	        return new SystemCapabilities(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.os = source["os"];
	        this.architecture = source["architecture"];
	        this.cpu_count = source["cpu_count"];
	        this.available_accelerators = source["available_accelerators"];
	        this.ffmpeg_version = source["ffmpeg_version"];
	        this.ffprobe_version = source["ffprobe_version"];
	        this.server_address = source["server_address"];
	        this.server_port = source["server_port"];
	    }
	}

}

