// @flow

import React from 'react'
import DropZone from 'react-dropzone'

const dropZoneStyle = {
  width: '100%',
  border: 'none'
}

export function RenderDropZone (onDropFuc: Function, fileName: string) {
  return (
    <DropZone
      accept={
        'application/pdf, \
        application/msword, \
        application/vnd.openxmlformats-officedocument.wordprocessingml.document, \
        application/vnd.oasis.opendocument.text, \
        application/x-zip-compressed, \
        application/x-gzip'
      }
      maxSize={5000000}
      onDrop={onDropFuc}
      style={dropZoneStyle}
      multiple={false}
      className={'file-upload'}
    >
      <i className='fas fa-paperclip hoverable' />
      {
        fileName ? <span>{fileName}</span> : <span>Select file...</span>
      }
    </DropZone>
  )
}
