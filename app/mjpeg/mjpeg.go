package mjpeg

/**
This  implements the saving jpeg frames to mp4 format buffer.
Supported formats:
  Uncompressed 8-bit (gray or indexed color), 24-bit (RGB),
  JPEG and PNG compression of individual frames
  16-bit and 32-bit (float) images are converted to 8-bit

* The AVI format written looks like this:
* RIFF AVI            RIFF HEADER, AVI CHUNK
*   | LIST hdrl       MAIN AVI HEADER
*   | | avih          AVI HEADER
*   | | LIST strl     STREAM LIST(s) (One per stream)
*   | | | strh        STREAM HEADER (Required after above; fourcc type is 'vids' for video stream)
*   | | | strf        STREAM FORMAT (for video: BitMapInfo; may also contain palette)
*   | | | strn        STREAM NAME
*   | | | indx        MAIN 'AVI 2.0' INDEX of 'ix00' indices
*   | LIST movi       MOVIE DATA (maximum approx. 0.95 GB)
*   | | 00db or 00dc  FRAME (b=uncompressed, c=compressed)
*   | | 00db or 00dc  FRAME
*   | | ...
*   | | ix00          AVI 2.0-style index of frames within this 'movi' list
* RIFF AVIX	          Only if required by size (this is AVI 2.0 extension)
*   | LIST movi       MOVIE DATA (maximum approx. 0.95 GB)
*   | | 00db or 00dc  FRAME
*   | | ...
*   | | ix00          AVI 2.0-style index of frames within this 'movi' list
* RIFF AVIX	          further chunks, each approx 0.95 GB (AVI 2.0)
* ...
**/
