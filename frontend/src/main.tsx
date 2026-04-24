import { StrictMode } from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { createRoot } from 'react-dom/client';
import '@/index.css';

import { Home } from '@/pages/Home.tsx';
import { AboutPage } from '@/pages/About.tsx';
import { Layout } from '@/components/layout/Layout.tsx';
import { History } from '@/pages/History.tsx';
import { Docs } from '@/pages/Docs.tsx';

import { TablePage } from '@/pages/TableRows.tsx';
import { RowForm } from '@/pages/RowForm.tsx';
import { TableForm } from '@/pages/TableForm.tsx';
import { NotFound } from '@/pages/NotFound.tsx';

import { MutationCache, QueryClient, QueryClientProvider } from '@tanstack/react-query';
import TableEditForm from '@/pages/TableEditForm.tsx';
import type { ErrorModel } from '@/client';
import { toast } from 'sonner';

const queryClient = new QueryClient({
	mutationCache: new MutationCache({
		onError: (error) => {
			const err = error as unknown as ErrorModel;
			toast.error(err.errors?.[0]?.message || err.detail || 'An unknown error occurred');
		}
	})
});

createRoot(document.getElementById('root')!).render(
	<StrictMode>
		<QueryClientProvider client={queryClient}>
			<BrowserRouter basename="/">
				<Routes>
					<Route element={<Layout />}>
						<Route path="/" element={<Home />} />

						<Route path="/new-table" element={<TableForm />} />
						<Route path="/tables/:tableName" element={<TablePage />} />
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
