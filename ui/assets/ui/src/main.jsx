// SPDX-FileCopyrightText: 2023 Iván Szkiba
//
// SPDX-License-Identifier: MIT

import React from 'react'
import ReactDOM from 'react-dom/client'
import { SSEProvider } from 'react-hooks-sse';
import { ThemeProvider, createTheme } from '@mui/material'
import App from './App'
import './index.css'

const base = new URLSearchParams(window.location.search).get("endpoint") || "http://localhost:5665/";

const theme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#7b65fa',
    },
    secondary: {
      main: '#A47D4F',
    },
  }
});

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <SSEProvider endpoint={base + "events"}>
      <ThemeProvider theme={theme}>
        <App />
      </ThemeProvider>
    </SSEProvider>
  </React.StrictMode>,
)