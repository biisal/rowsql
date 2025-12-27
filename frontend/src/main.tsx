import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import { Home } from "./pages/Home.tsx";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { AboutPage } from "./pages/about.tsx";
import { Layout } from "./components/Layout.tsx";
import { History } from "./pages/History.tsx";
import { Docs } from "./pages/Docs.tsx";

import { TableRows } from "./pages/TableRows.tsx";
import { RowForm } from "./pages/RowForm.tsx";
import { TabelForm } from "./pages/table-form.tsx";
import { NotFound } from "./pages/NotFound.tsx";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <BrowserRouter basename="/">
      <Routes>
        <Route element={<Layout />}>
          <Route path="/" element={<Home />} />
          <Route path="/table/:tableName" element={<TableRows />} />
          <Route path="/table/:tableName/form" element={<RowForm />} />
          <Route path="/about" element={<AboutPage />} />
          <Route path="/table/create/new" element={<TabelForm />} />
          <Route path="/history" element={<History />} />
          <Route path="/docs" element={<Docs />} />
          <Route path="*" element={<NotFound />} />
        </Route>
      </Routes>
    </BrowserRouter>
  </StrictMode>,
);
