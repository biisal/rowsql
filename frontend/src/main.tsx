import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import '@/index.css';

import {
	MutationCache,
	QueryClient,
	QueryClientProvider,
} from '@tanstack/react-query';
import { toast } from 'sonner';
import type { ErrorModel } from '@/client';
import { Layout } from '@/components/layout/Layout.tsx';
import { AboutPage } from '@/pages/About.tsx';
import { Docs } from '@/pages/Docs.tsx';
import { History } from '@/pages/History.tsx';
import { Home } from '@/pages/Home.tsx';
import { NotFound } from '@/pages/NotFound.tsx';
import { RowForm } from '@/pages/RowForm.tsx';
import TableEditForm from '@/pages/TableEditForm.tsx';
import { TableForm } from '@/pages/TableForm.tsx';
import { TablePage } from '@/pages/TableRows.tsx';
import { Editor } from '@/pages/Editor';

const queryClient = new QueryClient({
	mutationCache: new MutationCache({
		onError: (error) => {
			const err = error as unknown as ErrorModel;
			toast.error(
				err.errors?.[0]?.message || err.detail || 'An unknown error occurred',
			);
		},
	}),
});

createRoot(document.getElementById('root')!).render(
	<StrictMode>
		<QueryClientProvider client={queryClient}>
			<BrowserRouter basename="/">
				<Routes>
					<Route element={<Layout />}>
						<Route path="/" element={<Home />} />

						<Route path="/editor" element={<Editor />} />
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
