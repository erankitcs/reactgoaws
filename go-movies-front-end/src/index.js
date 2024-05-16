import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import { RouterProvider, createBrowserRouter } from 'react-router-dom';
import ErrorPage from './components/ErrorPage';
import Home from './components/Home';
import Movies from './components/Movies';
import Genres from './components/Genres';
import OneGenre from './components/OneGenre';
import EditMovie from './components/EditMovie';
import ManageCatalogue from './components/ManageCatalogue';
import GraphQL from './components/GraphQL';
import Login from './components/Login';
import Movie from './components/Movie';
import ManageMovieVideo from './components/ManageMovieVideo';
const router = createBrowserRouter(
  [
    {
      path: "/",
      element: <App />,
      errorElement: <ErrorPage />,
      children: [
        {
          index: true, element: <Home />
        },
        {
          path: "/movies", element: <Movies />
        },
        {
          path: "/movies/:id", element: <Movie />
        },
        {
          path: "/genres", element: <Genres />
        },
        {
          path: "/genres/:id", element: <OneGenre />
        },
        {
          path: "/admin/movie/0", element: <EditMovie />
        },
        {
          path: "/admin/movie/:id", element: <EditMovie />
        },
        {
          path: "/admin/movie/:id/upload", element: <ManageMovieVideo />
        },
        {
          path: "/manage-catalogue", element: <ManageCatalogue />
        },
        {
          path: "/graphql", element: <GraphQL />
        },
        {
          path: "/login", element: <Login />
        }
      ]
    }
  ]
)
const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    <RouterProvider router={router}></RouterProvider>
);
