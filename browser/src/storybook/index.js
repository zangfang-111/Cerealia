// @flow

import React from 'react'
import { action } from '@storybook/addon-actions'
import { storiesOf } from '@storybook/react'
import { Button } from '@storybook/react/demo'

storiesOf('Button', module)
  .add('with text', () => (
    <Button onClick={action('clicked')}>Hello Button</Button>
  ))
  .add('with some emoji', () => (
    <Button onClick={action('clicked')}><span role='img' aria-label='so cool'>😀 😎 👍 💯</span></Button>
  ))
