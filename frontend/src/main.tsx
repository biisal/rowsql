import { StrictMode } from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { createRoot } from 'react-dom/client';
import '@/index.css';

import { Home } from '@/pages/Home.tsx';
import { AboutPage } from '@/pages/about.tsx';
import { Layout } from '@/components/Layout.tsx';
import { History } from '@/pages/History.tsx';
import { Docs } from '@/pages/Docs.tsx';

import { TableRows } from '@/pages/TableRows.tsx';
import { RowForm } from '@/pages/RowForm.tsx';
import { TableForm } from '@/pages/table-form.tsx';
import { NotFound } from '@/pages/NotFound.tsx';

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import TableEditForm from '@/pages/TableEditForm.tsx';

const queryClient = new QueryClient();

createRoot(document.getElementById('root')!).render(
	<StrictMode>
		<QueryClientProvider client={queryClient}>
			<BrowserRouter basename="/">
				<Routes>
					<Route element={<Layout />}>
						<Route path="/" element={<Home />} />

						<Route path="/new-table" element={<TableForm />} />
						<Route path="/tables/:tableName" element={<TableRows />} />
						<Route path="/tables/:tableName/edit" element={<TableEditForm />} />
						<Route path="/tables/:tableName/rows/" element={<RowForm />} />

						<Route path="/about" element={<AboutPage />} />
						<Route path="/history" element={<History />} />
						<Route path="/docs" element={<Docs />} />
						<Route path="*" element={<NotFound />} />
					</Route>
				</Routes>
			</BrowserRouter>
		</QueryClientProvider>
	</StrictMode>,
);
