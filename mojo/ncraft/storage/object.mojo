
type Object {
    etag:          String    @1
    key:           String    @2 //< Name of the object
    last_modified: Timestamp @3 //< Date and time the object was last modified.
    size:          Int64     @4 //< Size in bytes of the object.
    content_type:  MediaType @5 //< A standard MIME type describing the format of the object data.
    content:       Bytes     @6
}
