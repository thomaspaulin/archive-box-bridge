# Copied from https://jacobtomlinson.dev/posts/2019/creating-github-actions-in-python/
FROM python:3.8-slim AS builder

RUN pip install pipenv

COPY Pipfile Pipfile.lock ./
RUN PIPENV_VENV_IN_PROJECT=1 pipenv install --deploy

# A distroless container image with Python and some basics like SSL certificates
# https://github.com/GoogleContainerTools/distroless
#FROM gcr.io/distroless/python3-debian10
FROM python:3.8-slim

COPY --from=builder /.venv /.venv
ENV PATH="/.venv/bin:$PATH"

WORKDIR /app
COPY link_finder.py link_renderer.py ./
#CMD ["./link_finder.py"]
ENTRYPOINT ["python", "link_finder.py"]