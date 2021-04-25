# Copied from https://jacobtomlinson.dev/posts/2019/creating-github-actions-in-python/
# Builder
FROM python:3.8-slim AS builder

RUN pip install pipenv

COPY Pipfile Pipfile.lock ./
RUN PIPENV_VENV_IN_PROJECT=1 pipenv install --deploy

# Runtime
FROM python:3.8-slim

COPY --from=builder /.venv /.venv
ENV PATH="/.venv/bin:$PATH"

WORKDIR /app
COPY link_finder.py link_renderer.py ./
ENTRYPOINT ["python", "/app/link_finder.py"]