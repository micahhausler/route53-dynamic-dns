FROM public.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base:latest

ARG TARGETARCH
ARG TARGETOS

#COPY route53-dynamic-dns /usr/bin/route53-dynamic-dns
COPY _output/bin/$TARGETOS-$TARGETARCH/route53-dynamic-dns /usr/bin/
ENTRYPOINT [ "/usr/bin/route53-dynamic-dns" ]