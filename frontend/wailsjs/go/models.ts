export namespace stylesage {
	
	export class Request {
	    text: string;
	    mode: string;
	    style: string;
	    context: string;
	
	    static createFrom(source: any = {}) {
	        return new Request(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.text = source["text"];
	        this.mode = source["mode"];
	        this.style = source["style"];
	        this.context = source["context"];
	    }
	}

}

